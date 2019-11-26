---
title: '[ubuntu]kubeadm安裝k8s集群(国内网络)'
categories:
  - 后端
  - 工具
tags:
  - kubernetes
  - kubeadm
comments: true
img: /img/kubernetes.png
abbrlink: 28805
date: 2019-11-23 13:31:43
---
## 步骤概览
> 1. 添加镜像源
> 2. 安装基础组件kubeadm、docker、kubelet
> 3. 初始化集群；
> 4. 补充网络组件；
> 5. 集群测试；

## 1. 添加镜像源
### 1.1 ubuntu 中添加镜像源,这里选用阿里云的源地址；
````
$ cat <<EOF > /etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt kubernetes-xenial main
EOF
$ sudo apt-get update
````

### 1.2 添加源之后，使用 apt-get update 命令会出现错误，原因是缺少相应的key，可以通过下面命令添加(BA07F4FB 为上面报错的key后8位)：

![](/illustration/ubuntu-kubernetes-install-error.png)

````
$gpg --keyserver keyserver.ubuntu.com --recv-keys 6A030B21BA07F4FB
$gpg --export --armor 6A030B21BA07F4FB | sudo apt-key add -
````

## 2. 安装基础组件kubeadm、docker、kubelet
````
$ sudo apt-get update && apt-get install -y docker.io kubelet  kubeadm
$ cat > /etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  },
  "storage-driver": "overlay2",
  "registry-mirrors": ["https://09icfnwb.mirror.aliyuncs.com"]
}
EOF
$ systemctl restart docker
$ docker version
````
Note: 这里配置的镜像仓库为官方镜像仓库，如果追求速度或者拉去镜像报错，可以参考配置使用[阿里镜像加速](https://cr.console.aliyun.com/cn-hangzhou/instances/mirrors)

## 3. 初始化集群
### 3.1 如果不关闭kubernetes运行会出现错误， 即使安装成功了，node重启后也会出现kubernetes server运行错误。
````
$sudo swapoff -a 
````
### 3.2 生成一个默认的kubeadm配置文件
````
$ mkdir ~/.kube/
$ kubeadm config print init-defaults > ~/.kube/config
````
生成的默认配置文件需要修改：
- `advertiseAddress` 原来是`1.2.3.4`,改为你的主机ip地址；
- `imageRepository` 原来是`k8s.gcr.io`,这里我们用阿里云的镜像源，所以需要改为`registry.cn-hangzhou.aliyuncs.com/google_containers`;
- `serviceSubnet` 原来是`10.96.0.0/12`,这里为了使用calico，我们需要改为`192.168.0.0/16`;
修改后的配置如下：
````yaml
apiVersion: kubeadm.k8s.io/v1beta2
bootstrapTokens:
- groups:
  - system:bootstrappers:kubeadm:default-node-token
  token: abcdef.0123456789abcdef
  ttl: 24h0m0s
  usages:
  - signing
  - authentication
kind: InitConfiguration
localAPIEndpoint:
  advertiseAddress: {YOUER IP ADDR}
  bindPort: 6443
nodeRegistration:
  criSocket: /var/run/dockershim.sock
  name: root
  taints:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
---
apiServer:
  timeoutForControlPlane: 4m0s
apiVersion: kubeadm.k8s.io/v1beta2
certificatesDir: /etc/kubernetes/pki
clusterName: kubernetes
controllerManager: {}
dns:
  type: CoreDNS
etcd:
  local:
    dataDir: /var/lib/etcd
imageRepository: registry.cn-hangzhou.aliyuncs.com/google_containers
kind: ClusterConfiguration
kubernetesVersion: v1.16.0
networking:
  dnsDomain: cluster.local
  serviceSubnet: 192.168.0.0/16
scheduler: {}
````
### 3.3 获取所有的镜像
````
$ kubeadm config images pull --config=${HOME}/.kube/config
````

### 3.4 初始化集群
````
$ kubeadm init --config=${HOME}/.kube/config
````
执行成功后按照提示配置config文件即可
### 3.5 创建完成后如果是单个节点的集群，需要去掉master节点的污点标记，使master节点也可以正常调度pod；
````
$ kubectl taint nodes --all node-role.kubernetes.io/master-
````

## 4. 补充网络组件

> [详情请参考官方文档](https://docs.projectcalico.org/v3.10/getting-started/kubernetes/)

````
$ kubectl apply -f https://docs.projectcalico.org/v3.10/manifests/calico.yaml
$ kubectl get pod -nkube-system -w
````
等待所有的pod都成Running状态，表示集群安装完成

## 5. 集群测试

````
$ kubectl run nginx --image=nginx:latest
$ kubectl expose deploy/nginx  --target-port=80 --port=8080
$ kubectl get pod |grep nginx
$ curl $(kubectl get svc|grep nginx |awk '{print$3}'):8080
````

可以正常访问nginx页面，表示集群创建成功，网络配置正常
````html
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
root@root:/etc/kubernetes/manifests#
root@root:/etc/kubernetes/manifests#
root@root:/etc/kubernetes/manifests# curl $(kubectl get svc|grep nginx |awk '{print$3}'):8080
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
````
