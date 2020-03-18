---
title: go module使用
categories:
  - 后端
tags:
  - go
comments: true
date: 2018-05-20 17:57:35
img:
---

### 外网问题
> https://www.jianshu.com/p/c5733da150c6 

### OLD
```
$ GO111MODULE=on GOPROXY=https://athens.azurefd.net go mod tidy
```

### New
```
$ GO111MODULE=on GOPROXY=https://mirrors.aliyun.com/goproxy/ go mod tidy
```
