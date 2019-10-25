---
title: k8s调度标签的使用
categories:
  - 后端
tags:
  - k8s
  - 调度
comments: false
abbrlink: 34402
date: 2019-10-23 22:38:30
img:
---

## 1. 如何管理master节点是否可以调度pod？
> 参考：
- https://blog.csdn.net/happyzwh/article/details/86063807
- https://cloud.tencent.com/developer/ask/148537 (添加label标签)


1. 去除master标签
```
$ kubectl taint node k8s-master node-role.kubernetes.io/master-

```

2. 去除master节点的污点标记；
```
$ kubectl lable node k8s-master node-role.kubernetes.io/master-

```

- 增加master标签，master节点不参与pod调度
```
$ kubectl taint nodes master1 node-role.kubernetes.io/master=:NoSchedule
```

## 2. 如何让pod调度到指定的节点上？
### 方法一 使用nodeSelector实现;
> https://kubernetes.io/docs/concepts/configuration/assign-pod-node/

Note: nodeSelector是一个硬性的选择器，只有当调度节点满足标签要求时才可调度成功
1. 给需要pod调度到的节点单上特殊标签
```
$ kubectl label node nodeName labelName=labelValue
```
2. 配置模板实例选择固定的label的节点调度
```
$ kubectl edit deployment youapp -nyouNamespace

nodeSelector:
  labelName: labelValue
```
### 方法二 使用affinity实现:


使pod调度到一个固定的节点上是一个双向选择的过程，需要完成两部分
1. 需要给node节点添加上污点，只能让特殊的pod调度上来；这是一个node选择pod的过程；
2. 需要给pod添加亲和性选择，让pod只能选择有固定标签的node节点，这是pod的一个选择过程；

```
affinity:
nodeAffinity:
  requiredDuringSchedulingIgnoredDuringExecution:
    nodeSelectorTerms:
    - matchExpressions:
      - key: fit
        operator: In
        values:
        - frm-systemdd
```