# Blog
Blog是基于hexo以及[hexo-theme-snippet]()主题的一个个人微博网站。

## 写文章
### 1. 生成文章主题
```
$ hexo new "【文章主题】我是一个demo"
INFO  =========================================
INFO    Welcome to use Snippet theme for hexo
INFO  =========================================
INFO  Created: F:\javascript\blog\source\_posts\【文章主题】我是一个demo.md
```
### 2. 编写文章内容
打开blog\source\_posts目录下的`【文章主题】我是一个demo.md`文件，编写你要发布的文章即可。
```
---
title: 【文章主题】我是一个demo
categories:
  - 前端
  - 分类2
tags:
  - 标签1
  - 标签2
comments: false
date: 2019-10-27 22:45:52
img:
---

# Demo
## 这是一个Demo演示
- 步骤一
- 步骤二
- 步骤三
```

## 发布
```
$ ./deploy.sh
_config.yml                                                                                                       100% 2221    58.2KB/s   00:00
go语言编码规范.md                                                                                           100% 6794   177.6KB/s   00:00
k8s调度标签的使用.md                                                                                       100% 1920    50.6KB/s   00:00
npm与node的安装与升级.md                                                                                   100% 1011    26.1KB/s   00:00
使用kind快速搭建一个k8s-ha集群.md                                                                       100% 2356    59.0KB/s   00:00
_config.yml                                                                                                       
deploy success !
```