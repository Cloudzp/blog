---
title: 使用kind快速搭建一个k8s-ha集群
categories:
  - 后端
tags:
  - kind
  - k8s
comments: false
img: 'https://github.com/kubernetes-sigs/kind/raw/master/logo/logo.png'
abbrlink: 34574
date: 2019-10-27 20:40:40
---
#
```yaml
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
```