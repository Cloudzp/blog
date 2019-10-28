---
title: 【git小知识】如何合并多个commit?
categories:
  - 后端
  - 工具
tags:
  - git
comments: true
abbrlink: 29953
date: 2019-10-28 22:12:04
img: /img/git.png
---

##  老铁的开发日志
那天风和日丽，老铁得到一个需求，要开发一个`csi`插件，老铁十分狂喜，终于可以大干一番了，于是啪啪啪，一顿编码，不到一下午，就写完了代码，本地验证没有问题，于是`git add ./*` `git commit -m "老铁NB"` `git push`一顿操作，将代码推送到了远端分支，提交了一个pr准备将代码合入主干分支，合入之前得让大佬进行代码`review`，大佬看完后啪啪啪一堆检视意见，老铁修改、提交 `git commit -m "老铁修改大佬检视意见"`，过一会又是啪啪啪几个检视意见，老铁再修改，提交 `git commit -m "老铁再次修改大佬检视意见"`，终于经过了两次修改，老铁改完了大佬的所有检视意见，最后询问大佬是否可以合入代码，这时有洁癖的大佬要求老铁将多个commit信息合并为一个，老铁心想怎样将多个commit合并为一个呢，老铁知道通过`git reset  commitID` 及`git push -f `可以进行commit合并，但是合并以后，大佬的检视意见就会没有了，这样操作老铁肯定会被大佬喷，老铁慌的一批，这时牛逼的同事出现了，"git rebase 可以合并commit而且能够保留检视意见"同事说到，下来就是给老铁一顿演示：

### 1.获取commitID 
找到你要合并的第一次的commit信息，获取到commitID
```
$ git log        
commit 5864c3ddb9e86221575aafdcaa95a0c972ba36b2
Author: Your Name <you@example.com>
Date:   Mon Oct 28 22:43:45 2019 +0800

    老铁再次修改大佬检视意见

commit 91a6f2a1d2d58fcde8228efb04b429e255049824
Author: Your Name <you@example.com>
Date:   Mon Oct 28 22:41:50 2019 +0800

    老铁修改大佬检视意见

commit 6727183434de98994fb9806ca4c5983a9f6e80a5
Author: Your Name <you@example.com>
Date:   Mon Oct 28 22:40:23 2019 +0800

    老铁NB

commit e1e1172fb677f92d493b3f3e6f2b531e90927aac

```

### 2. 通过rebase命令开启commit合并
通过`git rebase -i`命令打开合并窗口,这是你rebase后面的commit到你最后一次提交的commit信息都会出现在这里，而且是按照时间排序的，最上面一条是你提交的最早的一条commit，也是rebase命令后面的上一条commit，最下面一条是你最后提交的一条commit;要完成的就是将下面两条commit合并到第一条的commit上去，删除commit信息，保留修改信息；

```
$ git rebase -i e1e1172fb677f92d493b3f3e6f2b531e90927aac

pick 6727183 老铁NB
pick 91a6f2a 老铁修改大佬检视意见
pick 5864c3d 老铁再次修改大佬检视意见
....
```
### 3. 修改选择要合并的commit
编辑打开的编辑窗口，修改除了第一行以外的所有行的`pick`全部改为s(squash),然后`qw`保存，并退出：

```
pick 6727183 老铁NB
s 91a6f2a 老铁修改大佬检视意见
s 5864c3d 老铁再次修改大佬检视意见
```

保存后，stdout输出一下信息，表示合并成功；如果执行失败或是有冲突，请按照`5.知识扩展中的部署操作`
```
[detached HEAD a635fd9] 老铁NB
 Date: Mon Oct 28 22:40:23 2019 +0800
 1 file changed, 8 insertions(+)
Successfully rebased and updated refs/heads/develop.
```

### 4. 推送合并信息到远端分支
下来最后一步，push到远端分支，一定记住push需要加上`-f`参数，千万别按照提示，进行`git pull`否则你将陷入一个怪圈,
```
$ git push -f origin develop
```
push 成功以后，再次打开github，你会神奇的发现你的commit 已经合并成了一个，但是检视的信息还在，老铁太开心了！终于可以准点回家了。

### 5. 知识扩展：
#### a)如果在合并commit时候出现冲突，可以先解决冲突，然后执行`git add ./*`,之后执行`git rebase --continue`来完成冲突解决；
#### b)如果合并失败了，可以通过执行`git rebase --abort`来终止合并；
