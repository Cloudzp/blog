---
title: CGO
categories:
  - 后端
tags:
  - go
comments: true
abbrlink: 53379
date: 2018-06-22 19:57:35
img:
---

## 什么是CGO?
简单地说，cgo是在Go语言中使用C语言代码的一种方式。

## 为什么要有cgo?
C语言经过数十年发展，经久不衰，各个方面的开源代码、闭源库已经非常丰富。这无疑是一块巨大的宝藏，对于一门现代编程语言而言，如何用好现成的C代码就显得极为重要。

## 如何使用？
[见示例代码](main.go)

## 问题汇总？
1.golangd执行失败问题？
问题如下：
```
can't load package: package main: build constraints exclude all Go files in E:\workspace\src\FAQ\golang\src\cgo
```
> 解决方法参考：https://xs3c.co/archives/601

2.编译执行报错，错误如下：
```
$ CGO_ENABLED=1 go run main.go
# runtime/cgo
cc1.exe: sorry, unimplemented: 64-bit mode not compiled in
```
- 问题原因：
>解决编译器问题参考 https://github.com/mattn/go-sqlite3/issues/77

  1. 下载新的编辑器： http://tdm-gcc.tdragon.net/  （网络比较慢）；
  2. 安装完成后重新执行编译命令ok；
