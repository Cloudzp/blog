---
title: '在k8s通过pvc,pv的方式使用nfs'
categories:
  - 后端
tags:
  - k8s
  - nfs
  - pvc
comments: true
date: 2019-09-23 22:38:30
img:
---


## nfs存储的使用
nfs使用依赖nfs客户端,使用如下方式安装：
- 如果您使用CentOS、Redhat、Aliyun Linux操作系统，运行以下命令：

```
$ sudo yum install nfs-utils
```
- 如果您使用Ubuntu或Debian操作系统，运行以下命令：
```
$ sudo apt-get update
$ sudo apt-get install nfs-common
```

### 先创建pv,并填写文件存储信息
```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nas-test-pv
spec:
  capacity:
    storage: 5Gi
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: slow
  mountOptions:
    - vers=3
    - nolock
    - proto=tcp
    - noresvport
  nfs:
    path: /
    server: 0a6be494f8-fux86.cn-hangzhou.nas.aliyuncs.com
```
### 创建pvc 绑定pv
```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs-test-pvc
spec:
  accessModes:
    - ReadWriteOnce
  volumeMode: Filesystem
  resources:
    requests:
      storage: 5Gi
  storageClassName: slow
```

### 创建应用绑定pvc
```
ind: Deployment
metadata:
  labels:
    run: nginx
  name: nginx
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      run: nginx
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        run: nginx
    spec:
      containers:
      - image: nginx
        imagePullPolicy: Always
        name: nginx
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeDevices:
        - devicePath: /tmp2
          name: data
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - name: data
        persistentVolumeClaim:
          claimName: nfs-test-pvc
```

