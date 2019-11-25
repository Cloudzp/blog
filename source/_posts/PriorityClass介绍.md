---
title: Pod Priority and Preemption
categories:
  - kubernetes
tags:
  - k8s
comments: true
img: /img/kubernetes.png
abbrlink: 48493
date: 2019-11-19 22:10:34
---
## 1.简介
通过配置`PodPriority`可以为每个pod设置一个调度优先级（数字从0-2000000000），在集群资源不足时，可以防止pod被驱逐,在资源竞争过程中可以确保优先级高的pod优先被调度，`Preemption`是结核Pod Priority``特性一起使用，可以通过关闭`Preemption`,来防止资源不足时调度器会驱逐掉优先级底的pod；此特性在kubernetes1.11版本进入beta，在kubernetes1.14版本GA,从1.11版本以后就被默认打开；

版本的特性状态图
![](/illustration/priorityClass.png)

- 警告: 如果在一个不可信的kubernetes集群中使用此特性，一个恶意的用户可以通过创建一个最高级别的pod，导致其他用户的pod被驱逐或者不能被调度；要想解决这个问题可以参考[ResourceQuota](https://kubernetes.io/docs/concepts/policy/resource-quotas/)，ResourceQuota是对`PodPriority`的增强，管理员可以为特定优先级的用户创建`ResourceQuota`,来阻止创建优先级高的pod；

## 2. 如何使用`priority`和`preemption`
在kubernetes1.11及以后的版本使用`priority`和`preemption`特性，可以参考如下步骤：
- 1. 创建一个或多个`PriorityClasses`;

PriorityClasses yaml样例：
```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: high-priority
value: 1000000
globalDefault: false
description: "This priority class should be used for XYZ service pods only."
```

- 2. 创建一个pod使用`priorityClassName`去设置一个已经存在的`PriorityClass`,当然也可以不用直接创建pod，一般你可以设置`priorityClassName`
在Pod Template的对象集中，像 `Deployment`、`Daemonset`、`StatefulSet`、`Job`等。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
  labels:
    env: test
spec:
  containers:
  - name: nginx
    image: nginx
    imagePullPolicy: IfNotPresent
  priorityClassName: high-priority
```

如果你尝试该特性并且决定关闭它，你需要移除`PodPriority`的命令行参数，或者将它设置为false，之后重启`kube-apiserver`和`scheduler`。功能关闭之后，已经存在的pods保留他们的
`priority`属性，但是`preemption`是disable状态，以及`priority`属性是被忽略的。如果`PodPriority`已经关闭，你不能设置`priorityClassName`属性在新创建的pod中。

## 3. 如何关闭`preemption`(pod枪战：级别高的pod导致级别低的pod被驱逐，自己优先被调度)

> Note: 在kubernetes1.12+，当一个集群资源不足时，关键pod依赖`scheduler`抢占去调度。因此，不建议关闭`preemption`
> Note: 在kubernetes1.15及之后，如果`NonPreemptingPriority`特性是开启的，`PriorityClasses`有一个选项参数去设置`preemptionPolicy: Never`.
> 这将防止当前pods抢占其他的pod资源。

在Kubernetes1.11及之后版本，`preemption`是通过`kube-scheduler` flag `disablePreemption`控制的，默认时false，如果你想要的关闭preemption尽管提示不建议，
你可以通过设置`disablePreemption: true`.
这个选项仅对component配置中有效，在老版本的命令行参数是无效的，如下给出了component配置的关闭preemption的样例：
````yaml
apiVersion: kubescheduler.config.k8s.io/v1alpha1
kind: KubeSchedulerConfiguration
algorithmSource:
  provider: DefaultProvider

...

disablePreemption: true
````

## 4. PriorityClass
// TODO 