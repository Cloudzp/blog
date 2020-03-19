---
title: SNMP协议
categories:
  - 后端
tags:
  - snmp
comments: true
date: 2019-10-10 23:20:34
img:
---

> [可以学习一下](https://www.cnblogs.com/MYSQLZOUQI/p/5110900.html)
## 协议简介
   SNMP是专门ip网络设备设计的用来管理网络节点如：服务器，路由器，交换机等，是一种标准的应用层协议，基于tcp/ip;
   
### 1. 协议组成部分：
- SMI: 定义了SNMP框架所用信息的组织和标识，为MIB定义管理对象及使用管理对象提供模板.
- MIB: 定义了可以通过SNMP进行访问的管理对象的集合.一种树状数据库，MIB管理的对象，就是树的端节点，每个节点都有唯一位置和唯一名字.IETF规定管理信息库对象识别符（OID，Object Identifier）唯一指定，其命名规则就是父节点的名字作为子节点名字的前缀.
- SNMP协议: 络管理者如何对代理进程的MIB对象进行读写操作.

### 2. SNMP管理系统组成：
- 网络管理系统(一般是客户端)：　以该应用程序监视并控制被管理的设备，也称为管理实体（managingentity），网络管理员在这儿与网络设备进行交互。
- 被管理设备：　被管理的设备是一个网络节点，被管理的设备通过管理信息库（MIB）收集并存储管理信息，并且让网络管理系统能够通过SNMP代理者取得这项信息
- 代理者：代理者是一种存在于被管理的设备中的网络管理软件模块。代理者控制本地机器的管理信息，以和SNMP兼容的格式传送这项信息。　


## 测试环境搭建
### ubuntu上安装snmpd
```
$ sudo apt-get install snmpd
$ sudo service  snmpd start 
```
启动成功后ｓｎｍｐ会坚挺本地的`161`端口
```
$ netstat -nlp|grep 161 
…………
udp        0      0 127.0.0.1:161           0.0.0.0:*                           -
…………   
```

## 安装客户端工具：
```
$ sudo apt install -y net-snmp-utils
```
## go 代码中的应用
https://github.com/alouca/gosnmp 管理端lib
// TODO
