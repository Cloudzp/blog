---
title: 【mysql之】索引使用
categories:
  - 后端
tags:
  - mysql
  - db
comments: true
abbrlink: 23625
date: 2019-03-18 20:22:02
img:
---

### 索引（index）
- 优点：大大提高了数据库的检索效率；
- 缺点：索引其实也是一张表，该表存了索引与主键字段，并指向了实体表的记录，当进行INSERT 、
UPDATE、DELETE操作时不仅要修改实体表还要保存一下索引文件，所以如果过多的使用索引会使表的
修改操作变慢，还会浪费磁盘空间；

#### 添加索引
-  直接创建索引
````mysql
CREATE INDEX index_vpl_number2 ON t_park_record(vpl_number(32));

````

- 修改表结构添加索引
````mysql
ALTER TABLE t_park_record ADD INDEX index_vpl_number(vpl_number);

````
- 创建表时添加索引
````mysql

CREATE TABLE test_index(
  id INT NOT NULL PRIMARY KEY ,
  user_name VARCHAR(16) ,
  INDEX index_user_name (user_name(16))
);
````

### 删除索引
````mysql
DROP INDEX index_user_name ON test_index;
````
#### 修改表结构删除索引
````mysql
ALTER TABLE test_index2 DROP INDEX key_name_age ;

````

### 查询索引信息
```mysql
SHOW INDEX FROM  test_index2;
```

### 唯一索引

#### 创建唯一索引

````mysql
CREATE UNIQUE INDEX key_id_user_name ON test_index(`user_name`,`id`) ;
````

#### 修改表结构创建唯一索引
````mysql
ALTER TABLE test_index add UNIQUE key_id_user_name2(`user_name`,`id`); 
````

#### 创建表时创建唯一索引
````mysql
CREATE TABLE test_index2(
  id INT PRIMARY KEY NOT NULL ,
  name varchar(16),
  age INT ,
  UNIQUE key_name_age (`name`,`age`)
);
````

