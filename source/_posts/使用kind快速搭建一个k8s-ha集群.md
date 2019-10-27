---
title: 使用kind快速搭建一个k8s集群
categories:
  - 后端
tags:
  - kind
  - k8s
comments: false
img: https://github.com/kubernetes-sigs/kind/raw/master/logo/logo.png
abbrlink: 34574
date: 2019-10-27 20:40:40
---

## kind简介
Kind 是 Kubernetes In Docker 的缩写，顾名思义是使用 Docker 容器作为 Node 并将 Kubernetes 部署至其中的一个工具。官方文档中也把 Kind 作为一种本地集群搭建的工具进行推荐。

### 1. 安装kind
如果没有翻墙条件，需要通过编译源码的方式安装kind
```
$ GO111MODULE="on" GOPROXY=https://goproxy.io go get sigs.k8s.io/kind@v0.5.1
$ cp $GOPATH/bin/kind /usr/local/bin
$ kind version
v0.5.1
```

### 2. 创建集群
kind的node镜像需要翻墙下载，不能翻墙可以通过配置的方式获取

```
$ cat > kind.yaml << EFO
kind: Cluster
apiVersion: kind.sigs.k8s.io/v1alpha3
kubeadmConfigPatches:
- |
  apiVersion: kubeadm.k8s.io/v1beta1
  kind: ClusterConfiguration
  metadata:
    name: config
  networking:
    serviceSubnet: 10.0.0.0/16
  imageRepository: registry.aliyuncs.com/google_containers
  nodeRegistration:
    kubeletExtraArgs:
      pod-infra-container-image: registry.aliyuncs.com/google_containers/pause:3.1
- |
  apiVersion: kubeadm.k8s.io/v1beta1
  kind: InitConfiguration
  metadata:
    name: config
  networking:
    serviceSubnet: 10.0.0.0/16
  imageRepository: registry.aliyuncs.com/google_containers
nodes:
- role: control-plane
- role: worker
- role: worker
EFO
```

```
$ kind create cluster --name mycluster --config kind.yaml
Creating cluster "mycluster" ...
 ✓ Ensuring node image (kindest/node:v1.15.3) 🖼
 ✓ Preparing nodes 📦📦📦
 ✓ Creating kubeadm config 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
 ✓ Joining worker nodes 🚜
Cluster creation complete. You can now use the cluster with:

export KUBECONFIG="$(kind get kubeconfig-path --name="mycluster")"
kubectl cluster-info
```

### 3. 测试集群

```
$ export KUBECONFIG="$(kind get kubeconfig-path --name="mycluster")"
$ kubectl get node
NAME                      STATUS   ROLES    AGE     VERSION
mycluster-control-plane   Ready    master   3m29s   v1.15.3
mycluster-worker          Ready    <none>   2m51s   v1.15.3
mycluster-worker2         Ready    <none>   2m51s   v1.15.3
```