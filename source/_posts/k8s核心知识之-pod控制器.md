---
title: '[k8s核心知识之]　pod控制器'
categories:
  - 后端
tags:
  - k8s
comments: true
img: /img/kubernetes.png
abbrlink: 54325
date: 2020-03-28 12:07:54
---

Pod控制器由master的kube-controller-manager组件提供，常见的此类控制器有ReplicationController、ReplicaSet、Deployment、DaemonSet、
StatefulSet、Job和CronJob等，它们分别以不同的方式管理Pod资源对象

## 1. ReplicaSet 
### 1.1 功能简介
ReplicaSet（简称RS）是Pod控制器类型的一种实现，用于确保由其管控的Pod对象副本数在任一时刻都能精确满足期望的数量，ReplicaSet控制器资源启动后会查
找集群中匹配其标签选择器的Pod资源对象，当前活动对象的数量与期望的数量不吻合时，多则删除，少则通过Pod模板创建以补足，等Pod资源副本数量符合期望值后即
进入下一轮和解循环；

### 1.2 代码简介
```go
// 控制器结构体
type ReplicaSetController struct {
	// 资源类型及分组
	schema.GroupVersionKind

	kubeClient clientset.Interface
	podControl controller.PodControlInterface

    // 基于当前集群性能的一个最大的pod删除或创建数量　
	burstReplicas int
	// 核心方法，对pod数量进行定期维护，确保实际pod数量与rc期望中的一致
	syncHandler func(rsKey string) error
　　…………
}

// 此方法便是　syncHandler　的具体实现
// 主要做如下事情：
//  1. 根据key值从api中查找最新的rs对象；
//  2. 判断rs是否需要更新集群中的pod数量；
//  3. 列举出集群中的所有pod;
//  4. 通过labelselect过滤出所有rs关联的pod;
//  5. 根据实际pod数与rs中的replicas值决定删除或者创建pod;
//  6. 根据5中的结果更新rs的状态；
func (rsc *ReplicaSetController) syncReplicaSet(key string) error {

	startTime := time.Now()
	defer func() {
		klog.V(4).Infof("Finished syncing %v %q (%v)", rsc.Kind, key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	rs, err := rsc.rsLister.ReplicaSets(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(4).Infof("%v %v has been deleted", rsc.Kind, key)
		rsc.expectations.DeleteExpectations(key)
		return nil
	}
	if err != nil {
		return err
	}

	rsNeedsSync := rsc.expectations.SatisfiedExpectations(key)
	selector, err := metav1.LabelSelectorAsSelector(rs.Spec.Selector)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Error converting pod selector to selector: %v", err))
		return nil
	}

	// list all pods to include the pods that don't match the rs`s selector
	// anymore but has the stale controller ref.
	// TODO: Do the List and Filter in a single pass, or use an index.
	allPods, err := rsc.podLister.Pods(rs.Namespace).List(labels.Everything())
	if err != nil {
		return err
	}
	// Ignore inactive pods.
	filteredPods := controller.FilterActivePods(allPods)

	// NOTE: filteredPods are pointing to objects from cache - if you need to
	// modify them, you need to copy it first.
	filteredPods, err = rsc.claimPods(rs, selector, filteredPods)
	if err != nil {
		return err
	}

	var manageReplicasErr error
	if rsNeedsSync && rs.DeletionTimestamp == nil {
		manageReplicasErr = rsc.manageReplicas(filteredPods, rs)
	}
	rs = rs.DeepCopy()
	newStatus := calculateStatus(rs, filteredPods, manageReplicasErr)

	// Always updates status as pods come up or die.
	updatedRS, err := updateReplicaSetStatus(rsc.kubeClient.AppsV1().ReplicaSets(rs.Namespace), rs, newStatus)
	if err != nil {
		// Multiple things could lead to this update failing. Requeuing the replica set ensures
		// Returning an error causes a requeue without forcing a hotloop
		return err
	}
	// Resync the ReplicaSet after MinReadySeconds as a last line of defense to guard against clock-skew.
	if manageReplicasErr == nil && updatedRS.Spec.MinReadySeconds > 0 &&
		updatedRS.Status.ReadyReplicas == *(updatedRS.Spec.Replicas) &&
		updatedRS.Status.AvailableReplicas != *(updatedRS.Spec.Replicas) {
		rsc.enqueueReplicaSetAfter(updatedRS, time.Duration(updatedRS.Spec.MinReadySeconds)*time.Second)
	}
	return manageReplicasErr
}
```

### 1.2 小结：
对于rs来说并没有什么复杂的业务逻辑，其最核心的功能就是保持pod数与rs期望的replicas一致，多减少增。


## 2. Deployment控制器
### 2.1 功能简介
Deployment（简写为deploy）是Kubernetes控制器的又一种实现，它构建于ReplicaSet控制器之上，可为Pod和ReplicaSet资源提供声明式更新，相比较而言
，Pod和ReplicaSet是较低级别的资源，它们很少被直接使用，
Deployment控制器资源的主要职责同样是为了保证Pod资源的健康运行，其大部分功能均可通过调用ReplicaSet控制器来实现，同时还增添了部分特性:
- 事件和状态查看：必要时可以查看Deployment对象升级的详细进度和状态。
- 回滚：升级操作完成后发现问题时，支持使用回滚机制将应用返回到前一个或由用户指定的历史记录中的版本上。
- 版本记录：对Deployment对象的每一次操作都予以保存，以供后续可能执行的回滚操作使用。
- 暂停和启动：对于每一次升级，都能够随时暂停和启动。
- 多种自动更新方案：
  - Recreate: 即重建更新机制，全面停止、删除旧有的Pod后用新版本替代；
  - RollingUpdate: 即滚动升级机制，逐步替换旧有的Pod至新的版本;

Deployment控制器的滚动更新操作并非在同一个ReplicaSet控制器对象下删除并创建Pod资源，而是将它们分置于两个不同的控制器之下：旧控制器的Pod对象数量
不断减少的同时，新控制器的Pod对象数量不断增加，直到旧控制器不再拥有Pod对象，而新控制器的副本数量变得完全符合期望值为止，变动的方式和Pod对象的数量范
围将通过spec.strategy.rollingUpdate.maxSurge和spec.strategy.rollingUpdate.maxUnavailable两个属性协同进行定义。
- maxSurge：指定升级期间存在的总Pod对象数量最多可超出期望值的个数，其值可以是0或正整数，也可以是一个期望值的百分比；

###　2.2 代码简介
``` go
type DeploymentController struct {
	// rs　控制器，用来对rs进行增删改查操作
	rsControl     controller.RSControlInterface
	client        clientset.Interface
	eventRecorder record.EventRecorder

	// 核心方法　通过不断的根据key　找到对应的新的deploy对象，完成期望状态与现实状态的状态同步
	syncHandler func(dKey string) error
	// used for unit testing
	enqueueDeployment func(deployment *apps.Deployment)
　　　……
}

// 初始化控制器
func NewDeploymentController(dInformer appsinformers.DeploymentInformer, rsInformer appsinformers.ReplicaSetInformer, 
podInformer coreinformers.PodInformer, client clientset.Interface) (*DeploymentController, error) {
	……
	if client != nil && client.CoreV1().RESTClient().GetRateLimiter() != nil {
		if err := metrics.RegisterMetricAndTrackRateLimiterUsage("deployment_controller", client.CoreV1().RESTClient().
　　　　　GetRateLimiter()); err != nil {
			return nil, err
		}
	}
	dc := &DeploymentController{
		client:        client,
		eventRecorder: eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "deployment-controller"}),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "deployment"),
	}
	dc.rsControl = controller.RealRSControl{
		KubeClient: client,
		Recorder:   dc.eventRecorder,
	}

　　　// 每种控制器都会关注自己需要关注的资源对象，当资源对象发生改变的时候，通过钩子函数触发控制器进行某种操作；
     // 关注deployment资源的变化情况　
	dInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    dc.addDeployment,
		UpdateFunc: dc.updateDeployment,
		// This will enter the sync loop and no-op, because the deployment has been deleted from the store.
		DeleteFunc: dc.deleteDeployment,
	})
　　　
    //　关注rs资源的变化情况
	rsInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    dc.addReplicaSet,
		UpdateFunc: dc.updateReplicaSet,
		DeleteFunc: dc.deleteReplicaSet,
	})

    // 关注pod资源的变化情况
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: dc.deletePod,
	})

	dc.syncHandler = dc.syncDeployment
	dc.enqueueDeployment = dc.enqueue
　　　……
	return dc, nil
}

// 控制器的核心方法，用来持续同步deploy的状态，与实际状态保持一致
func (dc *DeploymentController) syncDeployment(key string) error {
	startTime := time.Now()
	klog.V(4).Infof("Started syncing deployment %q (%v)", key, startTime)
	defer func() {
		klog.V(4).Infof("Finished syncing deployment %q (%v)", key, time.Since(startTime))
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

　　 // 根据key值查找最新的Deployment对象　
	deployment, err := dc.dLister.Deployments(namespace).Get(name)
	if errors.IsNotFound(err) {
		klog.V(2).Infof("Deployment %v has been deleted", key)
		return nil
	}
	if err != nil {
		return err
	}

	// Deep-copy otherwise we are mutating our cache.
	// TODO: Deep-copy only when needed.
	d := deployment.DeepCopy()

	everything := metav1.LabelSelector{}
	if reflect.DeepEqual(d.Spec.Selector, &everything) {
		dc.eventRecorder.Eventf(d, v1.EventTypeWarning, "SelectingAll", "This deployment is selecting all pods. 
　　　　　A non-empty selector is required.")
		if d.Status.ObservedGeneration < d.Generation {
			d.Status.ObservedGeneration = d.Generation
			dc.client.AppsV1().Deployments(d.Namespace).UpdateStatus(d)
		}
		return nil
	}

　　 // 根据Deployment的label选取关联的所有rs资源
	rsList, err := dc.getReplicaSetsForDeployment(d)
	if err != nil {
		return err
	}
	// List all Pods owned by this Deployment, grouped by their ReplicaSet.
    // 列举当前Deployment所关联的所有pod资源，并按照rs进行分组；
	// podMap会被用于如下两种情况:
	// * 
	// * 检查Pod是否正确用pod-template-hash标签标记。
	// * 检查在重新创建部署的中间是否没有旧的Pod正在运行。
	podMap, err := dc.getPodMapForDeployment(d, rsList)
	if err != nil {
		return err
	}

	if d.DeletionTimestamp != nil {
		return dc.syncStatusOnly(d, rsList)
	}

	// Update deployment conditions with an Unknown condition when pausing/resuming
	// a deployment. In this way, we can be sure that we won't timeout when a user
	// resumes a Deployment with a set progressDeadlineSeconds.
	if err = dc.checkPausedConditions(d); err != nil {
		return err
	}

　　　// 如果滚动升级是终止状态则只进行pod水平伸缩，不能进行滚动升级或者回滚
	if d.Spec.Paused {
		return dc.sync(d, rsList)
	}

	// rollback is not re-entrant in case the underlying replica sets are updated with a new
	// revision so we should ensure that we won't proceed to update replica sets until we
	// make sure that the deployment has cleaned up its rollback spec in subsequent enqueues.
	if getRollbackTo(d) != nil {
		return dc.rollback(d, rsList)
	}

	scalingEvent, err := dc.isScalingEvent(d, rsList)
	if err != nil {
		return err
	}
	if scalingEvent {
		return dc.sync(d, rsList)
	}

　　// 根据Deployment的Strategy类型来进行相应的升级操作
	switch d.Spec.Strategy.Type {
	case apps.RecreateDeploymentStrategyType:
　　　　　// 进行Recreate策略
		return dc.rolloutRecreate(d, rsList, podMap)
	case apps.RollingUpdateDeploymentStrategyType:
        // 进行RollingUpdate策略
		return dc.rolloutRolling(d, rsList)
	}
	return fmt.Errorf("unexpected deployment strategy type: %s", d.Spec.Strategy.Type)
}

// Recreate策略　具体实现
// 这种策略的实现比较简单暴力，分为两步：
// 1. 遍历所有old　rs资源，将pod数都scale为0;
// 2. 创建一个新的rs，rs中的pod数就是Deployment中的replicas数量；
func (dc *DeploymentController) rolloutRecreate(d *apps.Deployment, rsList []*apps.ReplicaSet, 
　　　podMap map[types.UID][]*v1.Pod) error {
	
   // 创建一个新的rs，rs中的pod数就是Deployment中的replicas数量；
	newRS, oldRSs, err := dc.getAllReplicaSetsAndSyncRevision(d, rsList, false)
	if err != nil {
		return err
	}
	allRSs := append(oldRSs, newRS)
	activeOldRSs := controller.FilterActiveReplicaSets(oldRSs)

	// 遍历所有old　rs资源，将pod数都scale为0;
	scaledDown, err := dc.scaleDownOldReplicaSetsForRecreate(activeOldRSs, d)
	if err != nil {
		return err
	}
    …………
}

//RollingUpdate策略的具体逻辑，此中更新逻辑大致与Recreate策略一致，区别在于每一步更新RS的pod数量都是要经过严密的计算得到；
// 1. 创建一个新的rs资源，new rs中的replicas计算方式如下：
//    maxSurge如果是一个百分比会向上取证，也就是如果 deployment.Spec.Replicas＝10,maxSurge=25%,则maxSurge的int值为3;
//    (deployment.Spec.Replicas + maxSurge) - currentPodCount
// 2. scale　old　rs中的所有pod,具体缩减的pod数量需要根据　maxUnavailable数量来计算,大体计算方式如下：
//    maxUnavailable如果是一个百分比则会在计算过程中向下取证，eployment.Spec.Replicas＝10,maxUnavailable=25%,则maxUnavailable的int值为2;
//    minAvailable := deployment.Spec.Replicas - maxUnavailable
//    newRSUnavailablePodCount := newRS.Spec.Replicas - newRS.Status.AvailableReplicas
//    maxScaledDown := allPodsCount - minAvailable - newRSUnavailablePodCount
func (dc *DeploymentController) rolloutRolling(d *apps.Deployment, rsList []*apps.ReplicaSet) error {
	newRS, oldRSs, err := dc.getAllReplicaSetsAndSyncRevision(d, rsList, true)
	if err != nil {
		return err
	}
	allRSs := append(oldRSs, newRS)

	// Scale up, if we can.
	scaledUp, err := dc.reconcileNewReplicaSet(allRSs, newRS, d)
	if err != nil {
		return err
	}
	if scaledUp {
		// Update DeploymentStatus
		return dc.syncRolloutStatus(allRSs, newRS, d)
	}

	// Scale down, if we can.
	scaledDown, err := dc.reconcileOldReplicaSets(allRSs, controller.FilterActiveReplicaSets(oldRSs), newRS, d)
	if err != nil {
		return err
	}
	if scaledDown {
		// Update DeploymentStatus
		return dc.syncRolloutStatus(allRSs, newRS, d)
	}

	if deploymentutil.DeploymentComplete(d, &d.Status) {
		if err := dc.cleanupDeployment(oldRSs, d); err != nil {
			return err
		}
	}

	// Sync deployment status
	return dc.syncRolloutStatus(allRSs, newRS, d)
}

// 版本回滚的核心方法，回滚过程大体如下：
// 1.　在所有的pod
func (dc *DeploymentController) rollback(d *apps.Deployment, rsList []*apps.ReplicaSet) error {
	newRS, allOldRSs, err := dc.getAllReplicaSetsAndSyncRevision(d, rsList, true)
	if err != nil {
		return err
	}

	allRSs := append(allOldRSs, newRS)
	rollbackTo := getRollbackTo(d)
	// If rollback revision is 0, rollback to the last revision
	if rollbackTo.Revision == 0 {
		if rollbackTo.Revision = deploymentutil.LastRevision(allRSs); rollbackTo.Revision == 0 {
			// If we still can't find the last revision, gives up rollback
			dc.emitRollbackWarningEvent(d, deploymentutil.RollbackRevisionNotFound, "Unable to find last revision.")
			// Gives up rollback
			return dc.updateDeploymentAndClearRollbackTo(d)
		}
	}
	for _, rs := range allRSs {
		v, err := deploymentutil.Revision(rs)
		if err != nil {
			klog.V(4).Infof("Unable to extract revision from deployment's replica set %q: %v", rs.Name, err)
			continue
		}
		if v == rollbackTo.Revision {
			klog.V(4).Infof("Found replica set %q with desired revision %d", rs.Name, v)
			// rollback by copying podTemplate.Spec from the replica set
			// revision number will be incremented during the next getAllReplicaSetsAndSyncRevision call
			// no-op if the spec matches current deployment's podTemplate.Spec
			performedRollback, err := dc.rollbackToTemplate(d, rs)
			if performedRollback && err == nil {
				dc.emitRollbackNormalEvent(d, fmt.Sprintf("Rolled back deployment %q to revision %d", d.Name, 
　　　　　　　　　　rollbackTo.Revision))
			}
			return err
		}
	}
	dc.emitRollbackWarningEvent(d, deploymentutil.RollbackRevisionNotFound, "Unable to find the revision to rollback to.")
	// Gives up rollback
	return dc.updateDeploymentAndClearRollbackTo(d)
}


// rollback 回滚的核心逻辑，主要分为如下几步：
// 这里有一个很绕的逻辑，其实roolback并没对reversion及rs进行任何更新操作，它只是找到了roollbackto中的那个rs,然后将deploy中的template用rs完全替换，
// 在下一轮的更新中　DeploymentController　会按照一次正常的升级逻辑去处理这次回滚。
func (dc *DeploymentController) rollback(d *apps.Deployment, rsList []*apps.ReplicaSet) error {
　　　// 这里的只是为了拿到allOldRss
	newRS, allOldRSs, err := dc.getAllReplicaSetsAndSyncRevision(d, rsList, true)
	if err != nil {
		return err
	}

	allRSs := append(allOldRSs, newRS)
	rollbackTo := getRollbackTo(d)
	
    // 如果回滚的版本号为0,则回滚到最近版本后的一个版本 maxRevision-1
    // 如果要回滚到的版本号为0,则放弃回滚
	if rollbackTo.Revision == 0 {
		if rollbackTo.Revision = deploymentutil.LastRevision(allRSs); rollbackTo.Revision == 0 {
			// If we still can't find the last revision, gives up rollback
			dc.emitRollbackWarningEvent(d, deploymentutil.RollbackRevisionNotFound, "Unable to find last revision.")
			// Gives up rollback
			return dc.updateDeploymentAndClearRollbackTo(d)
		}
	}

    // 遍历所有的历史rs找到与rollbackTo的版本号匹配的rs
	for _, rs := range allRSs {
		v, err := deploymentutil.Revision(rs)
		if err != nil {
			klog.V(4).Infof("Unable to extract revision from deployment's replica set %q: %v", rs.Name, err)
			continue
		}
		if v == rollbackTo.Revision {
			klog.V(4).Infof("Found replica set %q with desired revision %d", rs.Name, v)
            // 1. 替换deployment的podTemplate为找到的对应rs的podTemplate 
            // 2. 去掉annotation中`deprecated.deployment.rollback.to` 头部
            // 3. 然后更新Deployment信息
　　　　　　　 // 在下一次getAllReplicaSetsAndSyncRevision调用期间，修订号将增加
			performedRollback, err := dc.rollbackToTemplate(d, rs)
			if performedRollback && err == nil {
				dc.emitRollbackNormalEvent(d, fmt.Sprintf("Rolled back deployment %q to revision %d", d.Name, 
　　　　　　　　　　rollbackTo.Revision))
			}
			return err
		}
	}
	dc.emitRollbackWarningEvent(d, deploymentutil.RollbackRevisionNotFound, "Unable to find the revision to rollback to.")
	// Gives up rollback
	return dc.updateDeploymentAndClearRollbackTo(d)
}
```

### 2.3 实操验证

#### 场景1. 我们对一个replicas=10的deployment进行滚动升级，并配置strategy为:Recreate
```
$ kubectl get pod -l run=nginx  -w
NAME                    READY   STATUS    RESTARTS   AGE
nginx-9ffc7d87b-4gxbr   1/1     Running   0          2m34s
nginx-9ffc7d87b-4mhvk   1/1     Running   0          2m34s
nginx-9ffc7d87b-6ftzn   1/1     Running   0          2m34s
nginx-9ffc7d87b-6zsrb   1/1     Running   0          2m34s
nginx-9ffc7d87b-8tw75   1/1     Running   0          2m34s
nginx-9ffc7d87b-p6jtw   1/1     Running   0          2m34s
nginx-9ffc7d87b-pxlcn   1/1     Running   0          2m34s
nginx-9ffc7d87b-q8wmb   1/1     Running   0          2m34s
nginx-9ffc7d87b-rsw79   1/1     Running   0          2m34s
nginx-9ffc7d87b-snkfv   1/1     Running   0          2m34s

$kubectl set image deploy/nginx nginx=nginx:1.0.0
………
nginx-9ffc7d87b-kjqqv   0/1     Terminating   0          6m36s
nginx-9ffc7d87b-kjqqv   0/1     Terminating   0          6m36s
nginx-9ffc7d87b-442fk   0/1     Terminating   0          6m36s
nginx-9ffc7d87b-442fk   0/1     Terminating   0          6m36s
nginx-9ffc7d87b-2ctmd   0/1     Terminating   0          6m36s
nginx-9ffc7d87b-2ctmd   0/1     Terminating   0          6m36s
nginx-9ffc7d87b-qmsds   0/1     Terminating   0          6m36s
nginx-9ffc7d87b-qmsds   0/1     Terminating   0          6m36s
nginx-9ffc7d87b-5gxfl   0/1     Terminating   0          6m37s
nginx-9ffc7d87b-5gxfl   0/1     Terminating   0          6m37s
nginx-9ffc7d87b-prxck   0/1     Terminating   0          6m37s
nginx-9ffc7d87b-prxck   0/1     Terminating   0          6m37s
nginx-f8b8788f8-qr528   0/1     Pending       0          0s
nginx-f8b8788f8-f6jcx   0/1     Pending       0          0s
nginx-f8b8788f8-qr528   0/1     Pending       0          0s
nginx-f8b8788f8-j6mjr   0/1     Pending       0          0s
nginx-f8b8788f8-j6mjr   0/1     Pending       0          0s
nginx-f8b8788f8-dm5wh   0/1     Pending       0          0s
…………
# 这里首先所有老的pod都先会被终止掉，所有pod终止完成后，才会创建新的pod
```

#### 场景2. 我们对一个replicas=10的deployment进行滚动升级，并配置strategy为:RollingUpdate,maxSurge:25%,maxUnavailable: 25%
```
$ kubectl get pod -l run=nginx  -w
NAME                    READY   STATUS    RESTARTS   AGE
nginx-9ffc7d87b-4gxbr   1/1     Running   0          2m34s
nginx-9ffc7d87b-4mhvk   1/1     Running   0          2m34s
nginx-9ffc7d87b-6ftzn   1/1     Running   0          2m34s
nginx-9ffc7d87b-6zsrb   1/1     Running   0          2m34s
nginx-9ffc7d87b-8tw75   1/1     Running   0          2m34s
nginx-9ffc7d87b-p6jtw   1/1     Running   0          2m34s
nginx-9ffc7d87b-pxlcn   1/1     Running   0          2m34s
nginx-9ffc7d87b-q8wmb   1/1     Running   0          2m34s
nginx-9ffc7d87b-rsw79   1/1     Running   0          2m34s
nginx-9ffc7d87b-snkfv   1/1     Running   0          2m34s

$kubectl set image deploy/nginx nginx=nginx:1.0.0

$ kubectl get pod -l run=nginx  -w                
NAME                    READY   STATUS        RESTARTS   AGE
nginx-9ffc7d87b-4gxbr   1/1     Running       0          3m38s
nginx-9ffc7d87b-4mhvk   1/1     Running       0          3m38s
nginx-9ffc7d87b-6ftzn   1/1     Running       0          3m38s
nginx-9ffc7d87b-6zsrb   1/1     Running       0          3m38s
nginx-9ffc7d87b-8tw75   1/1     Terminating   0          3m38s
nginx-9ffc7d87b-p6jtw   1/1     Running       0          3m38s
nginx-9ffc7d87b-pxlcn   1/1     Terminating   0          3m38s
nginx-9ffc7d87b-q8wmb   1/1     Running       0          3m38s
nginx-9ffc7d87b-rsw79   1/1     Running       0          3m38s
nginx-9ffc7d87b-snkfv   1/1     Running       0          3m38s
nginx-f8b8788f8-crlk7   0/1     Pending       0          2s   # 第一轮新创建
nginx-f8b8788f8-sfvm8   0/1     Pending       0          2s　　# 第一轮新创建
nginx-f8b8788f8-sq8qc   0/1     Pending       0          2s　　# 第一轮新创建
nginx-f8b8788f8-vttvf   0/1     Pending       0          1s　　# 第二轮新创建
nginx-f8b8788f8-vttvf   0/1     Pending       0          1s   # 第二轮新创建
nginx-f8b8788f8-dzk22   0/1     Pending       0          0s   # 第三轮新创建
nginx-f8b8788f8-dzk22   0/1     Pending       0          0s   # 第三轮新创建
# 执行完成升级命令后可以看到: 
根据滚动升级配置可知：　
maxSurge＝10 * 25%=3 这里百分比向上取整;
maxAvailable=10 * 25%=2 这里百分比向下取整;
- 第一轮创建：　创建了: (deployment.Spec.Replicas + maxSurge) - currentPodCount
                    =(10+3)-10=3 
- 第一轮删除：　终止了：　currentPodCount-(deployment.Spec.Replicas - maxUnavailable) -  newRSUnavailablePodCount 
                    =13-(10-2)-3=2

当2个old pod终止完成后　集群中的pod分布是：　3个new pod, 8个old pod= 11个pod
- 第二轮创建：　创建了:　(deployment.Spec.Replicas + maxSurge) - currentPodCount
                          = (10 + 10 * 25%) - 11= 2
- 第二轮删除：　终止了: currentPodCount-(deployment.Spec.Replicas - maxUnavailable) -  newRSUnavailablePodCount 
                     =13-(10-2)-3=2 (我们这里所有的新pod都没有启动成功); 

当2个old pod终止完成后　集群中的pod分布是：　5个new pod, 6个old pod= 11个pod ,根据上面的计算方法可以继续计算到会创建2个新的pod终止2个old
pod,一直循环直到所有的新pod被替换完成；　　　　　　　
```
###　2.4 小结：
deployment控制器是相当重要的一块内容，里面包含了所有无状态应用的升级更新逻辑，及整个升级过程中的pod新旧版本迭代变换，也是面试中很可能问到的内容，
主要记住如下重点即可：
#### 2.4.1 更新策略有两种：
- Recreate: 先终结掉所有旧版本的pod,等待所有pod全部终结后，重新创建；
- RollingUpdate:　滚动升级，需要关注两个值：
　- maxSurge: 可超出设置的replicas的最大数量,可以是int,也可以是百分比，如果是百分在计算过程中会向上取整数；
　- maxUnavailable: 最大的不可用pod数，同样可以配置为int,也可以配置成百分比，如果是百分比会在计算中向下取整；

#### 2.4.2 记住连个计算公式:
- 每次创建pod的数量计算公式：　deployment.Spec.Replicas+ maxSurge - currentPodCount
- 每次删除pod的数量计算公式：　currentPodCount - (deployment.Spec.Replicas-maxUnavailable) - newRsUnavailablePodCount

#### 2.4.3 回滚流程：
回滚流程中只是修改了Deployment中的podTemplate信息，具体的回滚操作是安装正常的升级策略完成的

#### 2.4.4 常用的命令：
```
# 替换镜像
$ kubectl set image deploy/{NAME} {CONTAINER_NAME}={NEW_IMAGE}
# 查看所有的历史版本 
$ kubectl rollout history deploy {NAME}
# 回滚操作
$ kubectl rollout undo  deploy {NAME}
# 查看更新或者回滚的状态
$ kubectl rollout status deploy {NAME}
# 暂停升级或者回滚操作
$ kubectl rollout pause deploy {NAME}
# 取消暂停操作
$ kubectl rollout resume deploy {NAME}
```
　　
