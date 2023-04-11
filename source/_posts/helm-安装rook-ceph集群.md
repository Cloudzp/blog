---
title: helm 安装rook-ceph集群
categories:
  - 后端
tags:
  - helm
  - ceph
comments: true
abbrlink: 11956
date: 2019-08-12 19:55:54
img:
---

# rook-ceph安装

## 安装步骤
- 下载最新的版本压缩包**rook-release-1.0.zip**
```
$ ll
rook-release-1.0.zip

$ unzip rook-release-1.0.zip
rook-release-1.0

$ cd rook-release-1.0/cluster/examples/kubernetes/ceph
$ kubectl apply -f common.yaml
$ kubectl apply -f operator.yaml
$ kubectl apply -f cluster.yaml // cluster.yaml文件按照如下修改配置
$ kubectl apply -f storageclass.yaml
```
- cluster.yaml
```yaml
apiVersion: ceph.rook.io/v1
kind: CephCluster
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"ceph.rook.io/v1","kind":"CephCluster","metadata":{"annotations":{},"name":"rook-ceph","namespace":"rook-ceph"},"spec":{"cephVersion":{"allowUnsupported":true,"image":"ceph/ceph:v14.2.1-20190430"},"dashboard":{"enabled":true},"dataDirHostPath":"/home/rook","mon":{"allowMultiplePerNode":true,"count":1},"network":{"hostNetwork":false},"rbdMirroring":{"workers":0},"storage":{"config":{"databaseSizeMB":"","journalSizeMB":"","osdsPerDevice":"1"},"deviceFilter":null,"directories":[{"path":"/home/rook"}],"useAllDevices":false,"useAllNodes":true}}}
  finalizers:
    - cephcluster.ceph.rook.io
  name: rook-ceph
  namespace: rook-ceph
spec:
  cephVersion:
    allowUnsupported: true
    image: ceph/ceph:v14.2.1-20190430
  dashboard:
    enabled: true
  dataDirHostPath: /home/rook
  mon:
    allowMultiplePerNode: true
    count: 1
    preferredCount: 0
  network:
    hostNetwork: false
  rbdMirroring:
    workers: 0
  storage:
    config:
      databaseSizeMB: ""
      journalSizeMB: ""
      osdsPerDevice: "1"
    directories:
      - config: null
        path: /home/rook
    useAllDevices: false
    useAllNodes: true
```

---

## 关于rook-ceph的问题汇总
   关于部署使用中遇到的问题可以优先查看官[issue](https://rook.io/docs/rook/v0.8/common-issues.html)是否有相关的问题：
 
### 问题1. docker重启后rook-osd启动不起来
```
rook-ceph-osd-0-576fc688c6-rs6hb      0/1     CrashLoopBackOff   10         26m
rook-ceph-osd-1-75d69db689-n6xp9      0/1     Error              10         26m
```
- 关注[issue](https://github.com/rook/rook/issues/3157)
- 临时解决方案，重新安装

### 问题2. 安装指导创建rook-ceph完成后无法创建rook-ceph-mgr 的pod
- 检查节点发现有残留的容器没有删除，删除残留容器后重新部署ok；
在每个节点上执行：
```
docker ps -a |grep rook
```

