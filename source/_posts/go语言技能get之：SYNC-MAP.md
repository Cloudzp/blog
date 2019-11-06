---
title: 【go语言技能get】：SYNC.MAP
categories:
  - 后端
tags:
  - go
comments: true
img: /img/go-logo.png
abbrlink: 997
date: 2019-11-06 21:14:19
---

# SYNC.MAP
sync.map是在1.9版本新增的功能，核心理念就是用空间换取时间，通过增加一个readmap来进行读取数据，从而避免使用锁来实现并发安全带来的性能消耗。

## 使用场景：
- 一次存储多次并发读取
- 多次并发修改已有的键值

## 实现原理
为了达到并发安全并且不损耗太多性能的前提下，sync.map通过原子操作来确保**查询**、**删除**操作的并发安全，存储操作也是通过加锁来实现的，所以sync.Map
相比于一般的map在大量的**读取、删除、更新**操作保证并发安全的同时，并没有太多的性能损失。

### 存储：
由于存储直接操作的是dirty对象，所以存储的过程需要加锁来完成；
- 判断当前的值是否存在read中，如果存在则直接更新read中的数据（更新是通过原子操作，直接操作read.m[key].entry.p的，所以不需要加锁完成）
- 如果不存在，则需要进行加锁操作(这里就是要添加新元素了)，由于加锁之前的操作是不受并发控制的，所以加锁后要重复去判断当前的值是否存在read中(也就是上一步的操作);这里会有三种情况：
  - 如果read中找到了key对应的entry，需要去判断当前的entry是否已经被删除掉了，如果被删除过，则将entry.p置为nil，并且将entry更新到dirty中，将数据存储在read[key].entry.p中(这步也是原子操作)
  - 如果read中找不到，但是dirty中找到了(这种数据一般是被删除掉的值)，直接将数据存储在read.m["key"].p中即可；
  - 如果read、dirty中都找不到，并且当前read相对与dirty并没有数据差异，则尝试从read复制dirty(这里是为了配合删除操作)，复制只会将未被标记为删除的数据copy到dirty中，并且标记read和dirty中存在数据差异，创建一个新的entry，并存储到dirty中，这里不会直接存储在read中因为read.m不能直接去添加数据；

### 查询：
- 先通过原子操作拿到read对象；
- 因为read.m是只读的，所以查询操作不用加锁即可完成查询，从read.m[key]中获取enry；
- 通过原子操作获取entry.p的值；
#### 有几种特殊情况
- 要查询的值


### 主要数据结构：
- Map
```go
type Map struct {
	// 同步锁 用来从dirty中读取写数据时使用；
	mu Mutex

    // read 使用原子操作实现的并发的写入，read的读取本身是并发安全的
	
	// Read 读取总是线程安全的， 但是存储数据时必须加上mu锁。
	
	// 存储在read中entries可以在不加锁的情况下更新，
	// 但是更新一个之前已经删除的entry时, entry 必须被copy到dirty 并且使用mu锁来控制并发安全。
	read atomic.Value // readOnly
	
	// dirty数据包含当前的map包含的entries,它包含最新的entries(包括read中未删除的数据,虽有冗余，但是提升dirty字段为read的时候非常快，
	// 不用一个一个的复制，而是直接将这个数据结构作为read字段的一部分),有些数据还可能没有移动到read字段中。
	
	// 对于dirty的操作需要加锁，因为对它的操作可能会有读写竞争。
	
	// 当dirty为空的时候， 比如初始化或者刚提升完，下一次的写操作会复制read字段中未删除的数据到这个数据中。
	dirty map[interface{}]*entry
	//
	// 当从Map中读取entry的时候，如果read中不包含这个entry,会尝试从dirty中读取，
	// 这个时候会将misses加一，
	// 当misses累积到 dirty的长度的时候， 就会将dirty提升为read,避免从dirty中miss太多次。
	// 因为操作dirty需要加锁。
	misses int
}
```
- readOnly
```go
// readOnly 被原子性的存储在Map.read字段中
type readOnly struct {
	m       map[interface{}]*entry
	// dirty中的map包含了某些key 但是在readOnly中却没有。
	amended bool // true if the dirty map contains some key not in m.
}
```
- expunged
expunged 是一个指针，标记着entry是否被从dirty的map中删除。

- entry
```go
// An entry is a slot in the map corresponding to a particular key.

// entry 是一个特殊key对应的一个插槽；
type entry struct {
	// 如果p==nil, entry已经被删除而且m.dirty==nil.
	// 
	// 如果p==expunged, 说明entry已经被删除，并且dirty已经被重新赋值过了，所以dirty中不会存在entry.
	//
	// 其他情况，entry是个正常值，如果m.dirty !=nil, entry存在于m.dirty中。
	//
	// 当m.dirty在下一次创建的时候，entry可能会被删除通过原子性的替换成nil,
	p unsafe.Pointer // *interface{}
}
```
### 主要函数
- Load
  
````go
// 返回当前key对应的value值
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	// 原子性操作从readOnly中获取read对象
	// 由于这里没有使用锁，所以读取的性能影响几乎为零。
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	// 如果不能从read中获取到key所对应的val值，并且dirty中存在read中的增量数据
	// 就通过锁的形式是去dirty中查找。
	if !ok && read.amended {
		m.mu.Lock()
		// 就在你判断的一瞬间，有可能dirty被升级，read已经被替换，所以这里需要在去read中获取一次；
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			// 无论entry是否存在，都要增加当前的miss次数。
			// 如果miss的次数大于等于dirty的长度，则升级dirty为read；
			// 升级后dirty将变为nil；我猜测这里是为了释放被删除掉的entry；
			// miss次数置为0；
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return nil, false
	}
	return e.load()
}
````
- Store

````go
// Store 存储一个键值对
// 可以看到存储过程都是在有锁的情况下进行的，所以存储时的并发效率影响挺大。
func (m *Map) Store(key, value interface{}) {
	// 和读取一样，先从read中尝试获取当前的值，如果获取到则直接返回。
	read, _ := m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}
    // 存储需要加锁
	m.mu.Lock()
	// 由于并发的原因，这里需要重复加索前的操作。
	read, _ = m.read.Load().(readOnly)
	if e, ok := read.m[key]; ok {
		// 确保entry没有被标记为擦除。
		// 如果entry没有被删除，则通过原子操作替换entry的p指针为nil。
		if e.unexpungeLocked() {
			// The entry was previously expunged, which implies that there is a
			// non-nil dirty map and this entry is not in it.
			m.dirty[key] = e
		}
		// 通过原子操作替换e.p的值。
		e.storeLocked(&value)
	// 如果read中找不到，就去dirty中获取
	// 如果获取到entry，则直接修改entry的值；
	} else if e, ok := m.dirty[key]; ok {
		e.storeLocked(&value)
	} else {
		// 如果read、dirty中都找不到，并且read中与dirty没有数据差异，
		// 说明这是一个初始化的过程？??
		if !read.amended {
			// We're adding the first new key to the dirty map.
			// Make sure it is allocated and mark the read-only map as incomplete.
			// 从read中重新创建dirty
			// 为什么这里要根据read创建dirty？ 什么时候dirty会是nil？
			// answer: 1. 每当dirty升级为read的时候dirty会被重置为nil，为什么这么做呢？
			// 在创建新元素的过程中，如果发现dirty为nil就根据read复制一个dirty，这里仅仅将没有清除的entry复制
			// 到dirty中。
			// 
			// answer：回答为什么dirty升级为read的过程中dirty会被置nil？
			// 删除过程中并不会直接去dirty中删除数据，而是在read中将entry的p指针置为nil，而等到dirty升级的过程中需要去清除掉dirtty中已经
			// 被删除掉的数据，所以dirty在升级完成后会被重置为nil，并且在新增元素的时候，dirty会从read中复制没有被标记为删除的数据；
			// TODO 那么为什么这么做呢？
		    // answer：提高删除的效率，删除是并不需要加索完成的。
			m.dirtyLocked()
			// 这里标记了read中与dirty中存在数据差异
			m.read.Store(readOnly{m: read.m, amended: true})
		}
		// 往dirty中添加数据
		m.dirty[key] = newEntry(value)
	}
	m.mu.Unlock()
}

````

- Delete

````go
// Delete 删除指定key的值
func (m *Map) Delete(key interface{}) {
	// 先从read中获取元素
	read, _ := m.read.Load().(readOnly)
	e, ok := read.m[key]
	// 如果获取不到，并且 dirty中存储在增量数据。
	// 也就是read中没有，可能dirty中有。
	
	// 这种情况一般是一个刚刚存储进来的数据，并且dirty还没有升级为read。
	if !ok && read.amended {
		m.mu.Lock()
		// 同理枷锁后需要重复枷锁前的操作。
		read, _ = m.read.Load().(readOnly)
		e, ok = read.m[key]
		if !ok && read.amended {
			// 从dirty的map中删除指定的key
			delete(m.dirty, key)
		}
		m.mu.Unlock()
	}
	// 如果read中有，则从read中擦除。
	// 一般的删除只会从read中删除，删除完成后，不会删除dirty中的数据，这里其实是用数据庸余的，
	// 在dirty升级的时候dirty会被清除为nil，这样来清除dirty中中的脏数据。
	if ok {
		// 删除read中的entry，即是：将entry中的p指针置为nil
		e.delete()
	}
}

// 删除一个entry
// 原子性的将entry的p指针替换为nil
// 注意这里仅仅修改了，read中的entry。
func (e *entry) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expunged {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

````
