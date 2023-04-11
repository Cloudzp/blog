---
title: linux 系统内核升级
categories:
  - 后端
tags:
  - linux
comments: true
abbrlink: 31402
date: 2019-03-22 23:11:57
img:
---

# 系统内核升级
## 添加可用的仓库
```
$ rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
$ rpm -Uvh https://www.elrepo.org/elrepo-release-7.0-3.el7.elrepo.noarch.rpm
```

## 查看可用的系统内核
```
$ yum --disablerepo="*" --enablerepo="elrepo-kernel" list available
```

## 安装最新的内核
```
$ yum --enablerepo=elrepo-kernel install kernel-ml
```

## 更新内核
```
$ grub2-mkconfig -o /boot/grub2/grub.cfg
```
## 查看系统内核
```
$ cat /boot/grub2/grub.cfg |grep menuentry
```
## 设置系统内核
```
$ grub2-set-default "CentOS Linux (5.0.3-1.el7.elrepo.x86_64) 7 (Core)"
```

