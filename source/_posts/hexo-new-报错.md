---
title: hexo new 报错
categories:
  - 前端
tags:
  - hexo
comments: true
date: 2018-09-18 21:11:18
img: 
---

前段时间换了电脑系统，重新下载了自己的hexo项目，在执行`hexo new "xxx"`过程中报了如下错误；

```
$ hexo new "【文章主题】我是一个demo"
ERROR Local hexo not found in ~/workdir/js/blog
ERROR Try running: 'npm install hexo --save'
```

## 解决方法
查找资料后，得知是因为node_modules文件夹的原因，需要删除重新安装

### 1.删除node_modules文件夹
```
$ rm -rf node_modules
```

### 2.重新安装依赖
```
$ cnpm install
```
### 3. 安装完成后再执行就不会报错了.
```
$ hexo new "【文章主题】我是一个demo"
INFO  =========================================
INFO    Welcome to use Snippet theme for hexo  
INFO  =========================================
INFO  Created: ~/workdir/js/blog/source/_posts/【文章主题】我是一个demo.md
```


