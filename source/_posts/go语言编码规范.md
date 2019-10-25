---
title: go语言编码规范
categories:
  - 后端
tags:
  - 编码指导
  - go
comments: false
abbrlink: 22572
date: 2019-10-23 22:31:26
img:
---
# GO 编码指导

[English](README.md) | [中文](README_ZH_CN.md)

> 语言官方编码规范指导 https://github.com/golang/go/wiki/CodeReviewComments#named-result-parameters
> 
> 参考ashleymcnamara的ppt分享

## 目标
> 通过此次交流分享会让你编写更容易阅读并且更容易维护的代码，想go团队一样
> - 怎样组织你的代码在pkg中以及这些pkg包含的内容。
> - 标准库中普遍存在的代码模式和约定;
> - 怎样去编写代码能让其更清晰容懂;
> - 不成文的Go会议超越了“go fmt”并让你看起来像一位资深Go贡献者写的 ；

## 大纲
> - 包管理
> - 命名约定
> - 源码约定

## 包管理
### a. 库的包管理
Go中有两个关键的代码组织区域会对代码的可用性，可测试性和功能产生巨大影响;
- Package 命名
- Package 管理

Packages 中应该只包含单一用途的代码，像：
- `archive`  `cmd` `crypto` `errors` `go` `index` `matl`

当一组package提供具有不同实现的通用功能集时，它们在父代下组织。
查看编码包的组织形式：
- `ascii85`  `base32` `binary` `encoding.go` `hex` 
- `asn1`     `base64` `csv`   `gob`         `json`

一些共性:
- Packages 名称描述他们的用途
- 很容易通过名称得知包的功能
- 名称通常很短
- 必要时，使用描述性父包和几个实现该功能的子项 - 如`encoding`包
 ```
  ├─encoding
  │  ├─charmap
  │  ├─htmlindex
  │  ├─ianaindex
  │  ├─internal
  │  │  └─identifier
  │  ├─japanese
  │  ├─korean
  │  ├─simplifiedchinese
  │  ├─testdata
  │  ├─traditionalchinese
  │  └─unicode
  │  ├─encoding.go // is the interface to encoding
  
 ```
我们见过的pkg都是库，它们旨在被某些可执行的程序引用并使用，例如：一个命令行工具。 

### b. 应用程序的包管理
你的应用程序的代码组织应该是什么样的？

当你有一个应用程序，与lib的代码组织略有不同。不同的是命令、可执行文件将所有这些包绑在一起。

应用程序包组织对系统的可测试性和功能性产生巨大影响。

在编写应用程序时，您的目标应该是编写易于理解，易于重构的代码容易让别人维持。

大多数libraries关注提供一个单个范围内的功能，例如：logging, ending, network access. 您的应用程序将所有这些库绑定在一起以创建工具或服务。 该工具或服务的范围将大得多。

在构建应用程序时，您应该将代码组织到包中，但这些包应该以两个类别为中心：
- Domain Types
- Services
`Domain Types` 是为业务功能和对象建模的类型。
`Services` 是在`Domain Types`上运行或与`Domain Types`一起运行的包。

[样例参考](https://medium.com/@benbjohnson/standaard-package-layout-7cdbc8391fc1)

包含`Domain Types`的包还应定义`Domain Types`与外部其他地方之间的接口。 这些接口定义了您要对`Domain Types`执行的操作。
- ProductService
- SupplierService
- AuthenticationService
- EmployeeStorage
- RoleStorage 

`Domain Types`包应该是应用程序存储库的根. 这使得打开代码库的任何人都清楚地知道正在使用什么类型，以及将对这些类型执行哪些操作。

`Domain Types`包或应用程序的根包不应具有任何外部依赖性。
- 它的存在是描述你的类型及其行为的唯一目的.

`Domain Types`接口的实现应该在独立的包中，按依赖性组织。 
依赖性包括:
- 外部数据源
- 传输逻辑 (http, RPC)

每个依赖项应该有一个包。

为什么每个依赖一个包？
- 更容易测试
- 易于替换/更换
- 没有循环依赖
              
## 命名约定
> 计算机科学有两个难点：缓存失效，命名和逐个错误
> Every developer on  Twitter 

命名很难，但是在你的类型，函数和包名中加入一些说明将使你的代码更具可读性。

### a. 包命名 
包名称应具有以下特征：
- 短
  - Prefer "transport" over "transportmechanisms"
- 清晰
  - 基于功能的清晰命名如:"bytes"
  - 用于描述外部依赖的实现的名称: "postgres"

软件包应该只为一个目的提供功能。 避免使用包如：
- util
- helpers
- etc.

他们的出现标志着你在哪里又缺少了一个`interface`，如：`util.ConvertOtherToThing()` 应该重构为Thinger的`interface`
这些包总是第一个遇到测试和循环依赖性问题的地方。 

### b. 变量命名 

变量名的一些常见约定：
- 使用camelCase而不是snake_case
- 使用单个字母变量来表示索引
  - `for i:=0; i < 10; i++ {}`
- 使用简短但描述性很强的变量名称
  - var count int
  - var cust Customer
Go中没有奖励点可以通过使用不必要的短变量来混淆代码。

使用变量的范围作为指南。 离您使用它的声明越远，名称应该越长。 

- 用重复的字母来表示一个`collection/slice/array`
  - `var tt []*Thing`
- 在`loop/range`内，使用单个字母
  - `for i, t := range tt{}`
  
这些约定在go的源代码中很常见。 

### c. 函数方法命名

避免使用重复包名称的包级别函数名称,包名已经声明了包的目的，因此不需要重复它。
- GOOD: `log.Info()`
- BAD: `log.LogInfo()`

Go代码没有setter和getter。
- GOOD: `custSvc.Customer()`
- BAD:  `custSvc.GetCustomer()`

如果你的`interface`只有一个方法，在方法名称后增加`-er`来作为`interface`的命名：
````go
type Stringer interface{
	String() string
}
````

如果你的`interface`中包含多个方法，使用一个能够概括这些方法功能的名称来命名：
````go
type CustomerStorage interface{
	Customer(id int)(*Customer, error)
	Save(c *Customer) error
	Delete(id int) error
}
````

在一个包内部将代码分成逻辑问题，如果包处理多种类型，请将每种类型的逻辑保留在其自己的源文件中：
````
package: postgres

orders.go
suppliers.go
products.go
````

## 源代码约定
在定义域对象的包中，为同一源文件中的每个对象定义类型和接口：
````
package: inventory

orders.go
-- contains Orders type and OrderStorage interface
````

总是用完整的句子做出说明。
````go
// An Order represents an order from a customer. 
type Order struct{}
````

使用`goimports`来管理你的导入，它们将始终按规范顺序排列。 标准库第一, 下来是外部库. 

避免使用else子句。 特别是在错误处理中。
````go
if err != nil {
	// error handling
	return  // or continue, etc.
}
// normal code
````