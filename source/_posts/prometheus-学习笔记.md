---
title: prometheus 学习笔记
categories:
  - monitor
tags:
  - 学习笔记
comments: true
abbrlink: 33511
date: 2022-04-21 22:55:47
img:
---
  
# 基本概念：
 - series：一个带有固定 label 的指标数据，一条变化的曲线
## 四种指标类型：
 Counter：非负递增的值，只增加不减
 Gauge：没啥限制，可增可减，适用于例如内存使用量这种指标。
 Histogram: 是桶统计类型，它自动将你定时的指标变成3个指标，例如你指定一个histogram叫a，那就会生成：
 - a_bucket{ le="某个桶值，例如100"} 用于统计打点值低于100的次数。
 - a_sum 表示打点值的和。
 - a_count 表示打点次数。
 Summary：百分位统计类型，会在客户端对于一段时间内（默认10分钟）的每个采样打点进行统计，并形成分位图。它也会生成3个指标，例如你指定一个Summary叫a，那就会生成
 - a{ quantile=0.99 } p99 百分位数，这个0.99是定义指标时指定的，可以指定好几个。
 - a_sum 表示打点值的和
 - a_count 表示打点次数