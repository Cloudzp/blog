---
title: ubuntu上npm与node的安装与升级
categories:
  - 前端
tags:
  - js
  - npm
  - node
comments: true
abbrlink: 599
date: 2019-10-25 21:42:35
img: /img/nodejs.png
---
## 安装nodejs
```
$ sudo apt-get install nodejs-legacy
```
这里安装的nodejs可能不是最新的，不过没关系，后面可以在线升级

## 安装npm
```
$ sudo apt-get install npm
```

## 升级nodejs(千万别先去升级npm，在升级nodejs之前！！！，可能会导致你的npm不可用)
```
$ sudo npm install -g n
$ sudo n stable   // 这里也可以切换到指定版本的nodejs e.g. sudo n v8.11.1
```
注意这里执行可能会遇到问题：
`cp: cannot stat `/usr/local/n/versions/node/7.10.0/lib': No such file or directory`
如果遇到如上问题，可通过卸载7.10.0版本的nodejs来解决 : `n - 7.10.0` 卸载后重新安装即可。

## 升级npm
```
$ sudo npm install npm@latest -g
```

## 安装cnpm, 一般在国内使用`npm`都比较慢，所以可以安装`cnpm`来使用阿里镜像源
```
$ npm install -g cnpm --registry=https://registry.npm.taobao.org
```

## 安装gulp
```
$ cnpm install --global gulp-cli
$ gulp --version
```