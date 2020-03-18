---
title: webservice
categories:
  - 后端
tags:
  - web
  - webservice
comments: true
date: 2019-05-18 20:22:18
img:
---

# webservice
一种古老的技术，至少对我来说是，主要目的是通过RPC构建一个分布式的可编程的web应用程序，各服务之间不用通过任何的第三方软件就可以相互通信获取数据；


Web Service平台需要一套协议来实现分布式应用程序的创建。任何平台都有它的数据表示方法和类型系统。要实现互操作性，Web Service平台必须提供一套标准的类型系统，
用于沟通不同平台、编程语言和组件模型中的不同类型系统。这些协议有：
- XML/XSD: xml是web Service的数据表示形式，XSD是web Service的数据类型系统。
- SOAP: (Sample Object Access Protocol)简单对象访问协议，用了web服务之间进行交换数据，是XML的一种简单封装。
- WSDL: Web Service的描述语言，基于XML格式来描述Web Service的函数、参数、返回值
- UDDI: UDDI 的目的是为电子商务建立标准；UDDI是一套基于Web的、分布式的、为Web Service提供的、信息注册中心的实现标准规范，同时也包含一组使企业能将自身提供的Web Service注册，以使别的企业能够发现的访问协议的实现标准。
