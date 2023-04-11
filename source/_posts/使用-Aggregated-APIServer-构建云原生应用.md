---
title: 使用 Aggregated APIServer 构建云原生应用
categories:
  - 后端
tags:
  - aggregated apiserver
  - kubernetes
comments: true
abbrlink: 5021
date: 2023-03-19 19:44:12
img:
---
# 1. Aggregated APIServer 是什么 ?
![img.png](img.png)

## 1.1 作用&功能 

## 1.2 APIServer 扩展的基本原理
kube-apiserver 作为整个 Kubernetes 集群操作 etcd 的唯一入口，负责 Kubernetes 各资源的认证&鉴权，校验以及 CRUD 等操作，提供 RESTful APIs，供其它组件调用：

![ 'image.png'](/download/attachments/1227785466/image-1637138768142.png?version=1&modificationDate=1637138768482&api=v2 'image.png')

**kube-apiserver其实包含三种APIServer：**
- **AggregatorServer**：负责处理 `apiregistration.k8s.io` 组下的APIService资源请求，同时将来自用户的请求拦截转发给 Aggregated APIServer(AA)；
- **KubeAPIServer**：负责对请求的一些通用处理，包括：认证、鉴权以及各个内建资源(pod, deployment，service)的 REST 服务等；
- **ApiExtensionsServer**：负责 CustomResourceDefinition（CRD）apiResources 以及 apiVersions 的注册，同时处理 CRD 以及相应 CustomResource（CR）的REST请求(如果对应 CR 不能被处理的话则会返回404)，也是 apiserver Delegation 的最后一环；

三个 APIServer 通过 delegation 的关系关联，在 kube-apiserver 初始化创建的过程中，首先创建的是 APIExtensionsServer，它的 delegationTarget 是一个空的 Delegate，即什么都不做，继而将 APIExtensionsServer 的 GenericAPIServer，作为 delegationTarget 传给了 KubeAPIServer，创建出了 KubeAPIServer，再然后，将 kubeAPIServer 的 GenericAPIServer 作为 delegationTarget 传给了 AggregatorServer，创建出了 AggregatorServer，所以他们之间 delegation 的关系为: Aggregator -> KubeAPIServer -> APIExtensions，如下图所示：
![ 'image.png'](/download/attachments/1227785466/image-1637138883298.png?version=1&modificationDate=1637138883612&api=v2 'image.png')


# 2. 为什么选择 Aggregated APIServer？
## 2.1 选择独立 API 还是 Aggregated APIServer ？
尽管使用 gin、go-restful 等 go 语言 web 框架可以轻易地构建出一个稳定的 API 接口服务，但以 kubernetes 原生的方式构建 API 接口服务还是有很多优势，例如:
- 能利用 kubernetes 原生的认证、授权、准入等机制，有更高的开发效率;
- 能更好的和 k8s 系统融合，借助 k8s 生态更快的推广自己的产品，方便用户上手;
- 借助于 k8s 成熟的 API 工具及规范，构建出的 API 接口更加规范整齐;

但是在很多场景下，我们还是不能确定到底使用聚合 API（Aggregated APIServer）还是独立 API 来构建我们的服务，官方为我们提供了两种选择的对比；如果你不能确定使用聚合 API 还是独立 API，下面的表格或许对你有帮助:

|  考虑 API 聚合的情况 |  优选独立 API 的情况 |
| ------------ | ------------ |
| 你在开发新的 API  | 你已经有一个提供 API 服务的程序并且工作良好  |
| 你希望可以是使用 `kubectl` 来读写你的新资源类别 |  不要求 `kubectl` 支持|
|你希望在 Kubernetes UI （如仪表板）中和其他内置类别一起查看你的新资源类别|不需要 Kubernetes UI 支持|
|你希望复用 [Kubernetes API 支持特性](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/api-extension/custom-resources/#common-features)|你不需要这类特性|
|你有意愿取接受 Kubernetes 对 REST 资源路径所作的格式限制，例如 API 组和名字空间。（参阅 [API 概述](https://kubernetes.io/zh/docs/concepts/overview/kubernetes-api/)）|你需要使用一些特殊的 REST 路径以便与已经定义的 REST API 保持兼容|
|你的 API 是[声明式的](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/api-extension/custom-resources/#declarative-apis)|你的 API 不符合[声明式](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/api-extension/custom-resources/#declarative-apis)模型|
|你的资源可以自然地界定为集群作用域或集群中某个名字空间作用域|集群作用域或名字空间作用域这种二分法很不合适；你需要对资源路径的细节进行控制|

**首先我们希望我们的 SKAI 平台能更好的和 k8s 结合，并且它是一个声明式的 API，尽可能的复用 Kubernets API 的特性，显然聚合 API 对我们来说更加适合。**

##  2.2 选择 CRDs 还是 Aggregated APIServer？
除了聚合 API，官方还提供了另一种方式以实现对标准 kubernetes API 接口的扩展：CRD（Custom Resource Definition ），能达到与聚合 API 基本一样的功能，而且更加易用，开发成本更小，但相较而言聚合 API 则更为灵活。针对这两种扩展方式如何选择，官方也提供了相应的参考。

通常，如果存在以下情况，CRD 可能更合适：
- 定制资源的字段不多；
- 你在组织内部使用该资源或者在一个小规模的开源项目中使用该资源，而不是在商业产品中使用；
  聚合 API 可提供更多的高级 API 特性，也可对其他特性进行定制；例如，对存储层进行定制、对 protobuf 协议支持、对 logs、patch 等操作支持。

两种方式的核心区别是定义 api-resource 的方式不同。在 Aggregated APIServer 方式中，api-resource 是通过代码向 API 注册资源类型，而 Custom Resource 是直接通过 yaml 文件向 API 注册资源类型。简单来说就是 CRD 是让 kube-apiserver 认识更多的对象类别（Kind），Aggregated APIServer 是构建自己的 APIServer 服务。虽然 CRD 更简单，但是缺少更多的灵活性，更详细的 CRDs 与 Aggregated API 的对比可参考[官方文档](https://kubernetes.io/zh/docs/concepts/extend-kubernetes/api-extension/custom-resources/#compare-ease-of-use)。

**对于我们而言，我们希望使用更多的高级 API 特性，例如 "logs" 或 "exec"，而不仅仅局限于 CRUD ，所以我们最终选择了 Aggregated APIServer 。**


# 3 如何开发一个 Aggregated APIServer ？
## 3.1 apiserver-build、sample-apiserver、 apiserver-runtime、 之间的关系。
- apiserver-build：是一个用于快速创建 Aggregated APIServer 的工具，它可以帮助我们快速创建项目骨架，并且使用 apiserver-builder 构建的项目目录结构比较清晰，更利于后期维护。
- sample-apiserver：是一个示例项目，它是一个简单的 Aggregated APIServer，我们可以参考它来实现自己的 Aggregated APIServer。
- apiserver-runtime：是一个用于构建 Aggregated APIServer 的库，它提供了一些基础的功能，例如：日志、配置、注册、认证、授权等，我们可以直接使用它来实现自己的 Aggregated APIServer。

apiserver-build 使用 apiserver-runtime 实现了 sample-apiserver 的功能，所以我们可以直接使用 apiserver-build 来创建项目骨架，然后再使用 apiserver-runtime 来实现自己的 Aggregated APIServer。

## 3.2 框架选型 apiserver-build & sample-apiserver
虽然官方提供了一个 [sample-apiserver](https://github.com/kubernetes/sample-apiserver)，我们可以参考实现自己的 Aggregated APIServer。但完全手工编写太过复杂，也不便于后期维护，我们最终选择了官方推荐的工具 [apiserver-builder](https://github.com/kubernetes-sigs/apiserver-builder-alpha)，apiserver-builder 可以帮助我们快速创建项目骨架，并且使用 apiserver-builder 构建的项目目录结构比较清晰，更利于后期维护。

### 3.2.1 apiserver-build 的使用

#### 3.2.1.1 安装 apiserver-builder 工具

通过 Go Get 安装
```
$ GO111MODULE=on go get sigs.k8s.io/apiserver-builder-alpha/cmd/apiserver-boot
```
通过安装包安装
- [下载](https://github.com/kubernetes-sigs/apiserver-builder-alpha/releases)最新版本
- 解压到 /usr/local/apiserver-builder/
- 如果此目录不存在，则创建此目录
- 添加/usr/local/apiserver-builder/bin到您的路径 `export PATH=$PATH:/usr/local/apiserver-builder/bin`
- 运行`apiserver-boot -h`

#### 3.2.1.2 初始化项目
完成 apiserver-boot 安装后，可通过如下命令来初始化一个 Aggregated APIServer  项目：
```
$ mkdir aa-demo
$ cd aa-demo 
$ apiserver-boot init repo --domain demo.io
```
执行后会生成如下目录：

```
.
├── BUILD.bazel
├── Dockerfile
├── Makefile
├── PROJECT
├── WORKSPACE
├── bin
├── cmd
│   ├── apiserver
│   │   └── main.go
│   └── manager
│       └── main.go -> ../../main.go
├── go.mod
├── hack
│   └── boilerplate.go.txt
├── main.go
└── pkg
    └── apis
        └── doc.go
```
- hack 目录存放自动脚本
- cmd/apiserver 是 aggregated server的启动入口
- cmd/manager 是 controller 的启动入口
- pkg/apis 存放 CR 相关的结构体定义，会在下一步自动生成

#### 3.2.1.3 生成自定义资源
```
$ apiserver-boot create group version resource --group animal --version v1alpha1 --kind Cat --non-namespaced=false
Create Resource [y/n]
y
Create Controller [y/n]
n
```
可根据自己的需求选择是否生成 Controller，我们这里暂时选择不生成, 对于需要通过 namespace 隔离的 resource 需要增加 --non-namespaced=false 的参数，默认都是 true。

执行完成后代码结构如下：
```
└── pkg
    └── apis
        ├── animal
        │   ├── doc.go
        │   └── v1alpha1
        │       ├── cat_types.go
        │       ├── doc.go
        │       └── register.go
        └── doc.go
```
可以看到在 pkg/apis 下生成了 animal 的 group 并在 v1alpha1 版本下新增了 `cat_types.go` 文件，此文件包含了我们资源的基础定义，我们在 spec 中增加字段定义，并在已经实现的 `Validate` 方法中完成基础字段的校验。
```
// Cat
// +k8s:openapi-gen=true
type Cat struct {
        metav1.TypeMeta   `json:",inline"`
        metav1.ObjectMeta `json:"metadata,omitempty"`

        Spec   CatSpec   `json:"spec,omitempty"`
        Status CatStatus `json:"status,omitempty"`
}

// CatSpec defines the desired state of Cat
type CatSpec struct {
	Name string `json:"name"`
}

func (in *Cat) Validate(ctx context.Context) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(in.Spec.Name) == 0 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "name"), in.Spec.Name, "must be specify"))
	}
	return allErrs
}
```

在 main 方法中进行资源注册。
```
func main() {
	err := builder.APIServer.
		WithResource(&animalv1alpha1.Cat{}).  // 资源注册。
		WithLocalDebugExtension().            // 本地调试。        
		DisableAuthorization().
		WithOptionsFns(func(options *builder.ServerOptions) *builder.ServerOptions {
			options.RecommendedOptions.CoreAPI = nil
			options.RecommendedOptions.Admission = nil
			return options
		}). // 控制参数， 不依赖 kube-apiserver 的 admission
		// WithoutEtcd(). // 不依赖 etcd 
		Execute()

	if err != nil {
		klog.Fatal(err)
	}
}

```

其他资源注册函数： 
- `WithResource`: 向 apiserver注册资源，使用默认的etcd备份存储。 
- `WithResourceAndStorage`: 向apiserver注册资源，创建一个新的etcd备份存储，GroupResource使用提供的策略。在大多数情况下，调用者应该使用WithResource 实现“apiserver-runtime/pkg/builder/rest”中定义的接口来控制策略。 
- `WithResourceAndHandler`: 为资源注册一个请求处理程序，而不是默认的 etcd备份的存储。
其他注册函数的使用样例参考：
// TODO



#### 3.2.1.4 部署运行
完成以上步骤，你其实已经拥有一个完整的 Aggregated APIServer，接下来我们试着将它运行起来；apiserver-boot 本身提供了两种运行模式：in-cluster、local; local 模式下只作为单独的 API 服务部署在本地方便做调试，过于简单这里不做过多介绍，主要关注一下 in-cluster 模式；in-cluster 可以将你的 Aggregated APIServer 部署在任何k8s集群中，例如：minikube，腾讯 TKE，EKS 等，我们这里使用 EKS 集群作为演示。

#### 4.1 创建[EKS集群](https://cloud.tencent.com/document/product/457/39813)&配置好本地[kubeconfig](https://cloud.tencent.com/document/product/457/39814)；
#### 4.2 执行部署命令 ；
```
$ apiserver-boot run in-cluster --image=xxx/demo.io/aa-demo:0.0.1 --name=aa-demo --namespace=default
```
在执行部署命令过程中，apiserver-boot 主要帮我们做了如下几件事情：
- 自动生成 APIServer Dockerfile 文件；
- 通过 APIServer Dockerfile 构建服务镜像，并将镜像推送到指定仓库；
- 在config目录下生成 CA 及其他 APIServer 部署需要的证书文件；
- 在config目录下生成 APIServer 部署需要的 Deployment、Service、APIService、ServiceAccount 等 yaml 文件；
- 将上一步生成的 yaml 文件部署到集群中；

#### 5.功能验证
### 5.1 确认 Resource 注册成功；
```
$ kubectl api-versions |grep animal
animal.demo.io/v1alpha1
```

### 5.2 确认 Aggregated APIServer 能正常工作；
```
$ kubectl get apiservice v1alpha1.animal.demo.io 
NAME                      SERVICE             AVAILABLE   AGE
v1alpha1.animal.demo.io   default/demo   True        19h
```

### 5.3 创建并查看新增的 Resource
创建
```
$ cat lucky.yaml
apiVersion: animal.skai.io/v1alpha1
kind: Cat
metadata:
  name: mycat
  namespace: default
spec:
  name: lucky
  
# 创建自定义 resource
$ kubectl apply -f lucky.yaml
```
查找
```
# 查找自定义 resource 列表
$ kubectl get cat
NAME    CREATED AT
mycat   2021-11-17T09:08:10Z

# 查找自定义资源详情
$ kubectl get cat mycat -oyaml
apiVersion: animal.skai.io/v1alpha1
kind: Cat
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"animal.skai.io/v1alpha1","kind":"Cat","metadata":{"annotations":{},"name":"mycat"},"spec":{"name":"lucky"}}
  creationTimestamp: "2021-11-17T09:08:10Z"
  name: mycat
  resourceVersion: "17"
  uid: 98af0905-f01d-4042-bad3-71b96c0919f4
spec:
  name: lucky
status: {}
```

### 3.2.2 sample-apiserver 的使用
sample-apiserver 是 Kubernetes 官方提供的一个 Aggregated APIServer 的样例，它提供了一个简单的 API，可以用来学习和测试 Aggregated APIServer 的功能，我们可以通过如下命令快速部署一个 sample-apiserver：
```
staging/src/k8s.io/sample-apiserver
├── artifacts
│   ├── example
│   │   ├── apiservice.yaml
      ...
├── hack
├── main.go
└── pkg
├── admission
├── apis
├── apiserver
├── cmd
├── generated
│   ├── clientset
│   │   └── versioned
              ...
│   │       └── typed
│   │           └── wardle
│   │               ├── v1alpha1
│   │               └── v1beta1
│   ├── informers
│   │   └── externalversions
│   │       └── wardle
│   │           ├── v1alpha1
│   │           └── v1beta1
│   ├── listers
│   │   └── wardle
│   │       ├── v1alpha1
│   │       └── v1beta1
└── registry
```
其中，artifacts用于部署yaml示例
- hack目录存放自动脚本(eg: update-codegen)
- main.go是aggregated server启动入口；pkg/cmd负责启动aggregated server具体逻辑；pkg/apiserver用于aggregated server初始化以及路由注册
- pkg/apis负责相关CR的结构体定义，自动生成(update-codegen)
- pkg/admission负责准入的相关代码
- pkg/generated负责生成访问CR的clientset，informers，以及listers
- pkg/registry目录负责CR相关的RESTStorage实现

### registry 的实现
```

```


### 部署
在开发过程中，独立运行样本 apiserver 很有帮助，即没有 一个 Kubernetes API 服务器，用于 authn/authz，没有聚合。这是可能的，但需要 几个标志、密钥和证书如下所述。你仍然需要一些 kubeconfig， 例如，但 Kubernetes 集群不用于 authn/z。迷你库贝或 hack/本地 cluster.sh 集群将正常工作。~/.kube/config
所描述的设置不是信任 kube-apiserver 中的聚合器，而是使用本地 基于客户端证书的 X.509 身份验证和授权。这意味着客户端 证书受 CA 信任，并且传递的证书包含组成员身份 到组。当我们使用 禁用委派授权时， 只有此超级用户组获得授权。system:masters--authorization-skip-lookup

首先，我们需要一个 CA 来签署客户端证书：
```
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt
```
然后，我们为超级用户组中的用户创建由此 CA 签名的客户端证书：developmentsystem:masters
```
openssl req -out client.csr -new -newkey rsa:4096 -nodes -keyout client.key -subj "/CN=development/O=system:masters"
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -sha256 -out client.crt
```
由于 curl 需要带有密码的 p12 格式的客户端证书，因此请进行转换：
```
openssl pkcs12 -export -in ./client.crt -inkey ./client.key -out client.p12 -passout pass:password
```
使用这些密钥和证书后，我们启动服务器：

```
etcd &
sample-apiserver --secure-port 8443 --etcd-servers http://127.0.0.1:2379 --v=7 \
   --client-ca-file ca.crt \
   --kubeconfig ~/.kube/config \
   --authentication-kubeconfig ~/.kube/config \
   --authorization-kubeconfig ~/.kube/config
```

第一个 kubeconfig 用于共享告密者访问 Kubernetes 资源。传递给的第二个 kubeconfig 用于满足委托 身份验证器。传递给的第三个 kubeconfig 用于满足委托 授权方。身份验证者和授权者都不会 实际使用：由于，我们的开发X.509 证书被接受并验证我们作为会员的身份。 是超级用户组，以便委派 跳过授权。--authentication-kubeconfig--authorized-kubeconfig--client-ca-filesystem:masterssystem:masters

使用 curl 使用 p12 格式的客户端证书访问服务器进行身份验证：
```
curl -fv -k --cert-type P12 --cert client.p12:password \
   https://localhost:8443/apis/wardle.example.com/v1alpha1/namespaces/default/flunders
```   
或者使用 wget：

```
wget -O- --no-check-certificate \
   --certificate client.crt --private-key client.key \
   https://localhost:8443/apis/wardle.example.com/v1alpha1/namespaces/default/flunders
```   
注意：最近的 OSX 版本使用 curl 破坏了客户端证书。在 Mac 上尝试，然后：brew install httpie

```
http --verify=no --cert client.crt --cert-key client.key \
   https://localhost:8443/apis/wardle.example.com/v1alpha1/namespaces/default/flunders
``` 

# 总结：
本文从实战角度出发介绍我们开发 SKAI 平台过程中选择 Aggregated API 的原因，以及 kube-apisever 的扩展原理，最后介绍了 apiserver-builder 工具，并演示如何一步一步构建起自己的 Aggregated API，并将它部署到 EKS 集群中。希望该篇 Aggregated APIServer 最佳实践可以帮助即将使用 k8s API 扩展来构建云原生应用的开发者。