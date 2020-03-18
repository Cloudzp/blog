---
title: helm工具的安装
categories:
  - 后端
  - 工具
tags:
  - helm
  - k8s
comments: true
date: 2019-07-10 20:51:54
img:
---

> https://github.com/helm/helm/releases/tag/v2.14.1
## 1.下载安装包解压
- [releases](https://github.com/helm/helm/releases/tag/v2.14.1)

```
$ tar -xzvf helm-v2.14.1-linux-amd64.tar.gz.asc
$ cp  helm-v2.14.1-linux-amd64/bin/helm /usr/local/bin
```

## 2. 安装
```
# 初始化安装
helm init

# Create a ServiceAccount for Tiller in the `kube-system` namespace
kubectl --namespace kube-system create sa tiller

# Create a ClusterRoleBinding for Tiller
kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller

# Patch Tiller's Deployment to use the new ServiceAccount
 {}kubectl --namespace kube-system patch deploy/tiller-deploy -p '{"spec": {"template": {"spec": {"serviceAccountName": "tiller"}}}}'

```

