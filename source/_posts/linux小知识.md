---
title: linux小知识
categories:
  - 后端
tags:
  - linux
comments: true
img: /img/linux.jpg
abbrlink: 55295
date: 2019-11-19 23:30:42
---
# linux的常用技巧总结
## 1. ubuntu设置sudo免密：
```
$ sudo visudo

```
在最后一行增加如下配置,USER_NAME使用你要配置的用户替换；
```markdown
${USER_NAME} ALL=(ALL) NOPASSWD: ALL

```

![](/illustration/ubuntu_sudo.png)


## 2. ubuntu替换apt镜像源为阿里的镜像源；
系统默认带的镜像源都是国外源，通过`apt`下载软件比较慢，有的还有可能需要翻墙，所以我们可以通过替换为阿里源的方式解决，方法如下：

### 2.1 方法一：
```
$ sudo mv /etc/apt/sources.list /etc/apt/sourses.list.backup
$ sudo cat >  /etc/apt/sources.list << EFO
deb http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
EFO
```

### 2.2 方法二：
```
$ sudo mv /etc/apt/sources.list /etc/apt/sourses.list.backup
$ sudo vi /etc/apt/sources.list 
deb http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
```
vi 打开文件后输入`:0,$s/archive.ubuntu.com/mirrors.aliyun.com/g`

### 2.3 方法三:
```
$ sudo mv /etc/apt/sources.list /etc/apt/sourses.list.backup
$ sudo sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
```

### 2.4 几个知名国内源地址

````
#中科大源
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse

#阿里云源
deb http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse

#清华源
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse
````



## 3. ssh安装

> https://help.ubuntu.com/lts/serverguide/openssh-server.html (官方指导)

```
$ sudo apt install openssh-server
$ sudo cp /etc/ssh/sshd_config /etc/ssh/sshd_config.original
$ sudo chmod a-w /etc/ssh/sshd_config.original
$ sudo systemctl restart sshd.service
```

## 4. 打造一个高逼格的终端
工欲善其事必先利其器，一个高逼格的终端，不但可以提高我们的工作效率，更能让我们更有面，心情好了，干什么事情都有激情；废话不说先上图一张：
- 分屏随意分
![](/illustration/deepin_term.png)
![](/illustration/deepin_term_split.png)
- 文件随意打开
![](/illustration/deepin_term_openfile.png)
- 主题随意换
![](/illustration/deepin_term_theme.png)
- 远程链接管理
![](/illustration/deepin_term_sshm.png)

这就是今天的猪(主)脚(角)`深度终端`，听这名字都很有深度是吧，再看看看它高逼格的外表，骚骚的操作，真是内外兼修骚气冲天，你不想来一发？要来？那么请跟随我的步伐。 
### 4.1 ubuntu上安装`深度终端`，如果你使用的是`deepin`操作系统，那么打扰了，你可以告辞了；
#### 步骤 1. 更新apt源
````
$ sudo cat > /etc/apt/sources.list.d/deepin-ubuntu-dde-bionic.list << EOF
deb https://mirrors.aliyun.com/deepin/ panda main contrib non-free
deb-src https://mirrors.aliyun.com/deepin/ panda main contrib non-free
EOF
$ sudo apt update
````
添加mirror以后使用update遇到错误
`The following signatures couldn't be verified because the public key is not available: NO_PUBKEY 625956CB3E21DE5`
修复方式如下：
````
$ sudo wget https://mirrors.aliyun.com/deepin/project/deepin-keyring.gpg
$ sudo gpg --import deepin-keyring.gpg
$ sudo gpg --export --armor B3E21DE5 | sudo apt-key add -    // B3E21DE5 就是NO_PUBKEY的后8位
$ sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 625956CB3E21DE5

````
####  步骤 2. 安装深度终端
````
$ sudo apt update && sudo apt install -y deepin-terminal
````
**到这里我们的男猪脚`深度终端`基本已经表演完了，下来该我们的伴娘（oh-my-zsh）出场了**

### 4.2 oh-my-zsh安装

#### 步骤 1. 安装zsh，由于oh-my-zsh依赖，zsh所以安装之前需要先行安装zsh
````
$ sudo apt install zsh
````
or
````
$ sudo yum install zsh
````
#### 步骤 2. 安装oh-my-zsh
- via curl
````
$ sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
````

- via wget
````
$ sh -c "$(wget -O- https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
````
#### 步骤 3. 配置主题
主题可根据个人癖好选择，我比较喜欢用`ys`，这个主题风格比较骚，符合我的气质，广大码农可根据自己的癖好选择[主题传送门](https://github.com/ohmyzsh/ohmyzsh/tree/master/themes)
打开`~/.zshrc`文件，搜索`ZSH_THEME`，修改如下配置`ZSH_THEME="ys"`;

#### 步骤 4. 配置插件
个人觉得一个比较实用的插件`zsh-autosuggestion`,是我比较喜欢的，可以帮助我们补充命令，提高效率；
````
$ git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions

````
打开`~/.zshrc`文件，搜索`plugins`，增加如下配置`plugins=(git zsh-autosuggestions)`;

## 5. ssh免密
````
$ ssh-keygen -t rsa
Generating public/private rsa key pair.
Enter file in which to save the key (/home/cloudzp/.ssh/id_rsa): 
Enter passphrase (empty for no passphrase): 
Enter same passphrase again: 
Your identification has been saved in /home/cloudzp/.ssh/id_rsa.
Your public key has been saved in /home/cloudzp/.ssh/id_rsa.pub.
The key fingerprint is:
SHA256:8LQHvXSlaYz0Ut2MgcskVud8xIshP8isjT1jrtCR7PI cloudzp@deepin
The key's randomart image is:
+---[RSA 2048]----+
|          ..oo+*.|
|         oo*oOo.+|
|      . o.B+@.= o|
|       +.+.Oo+ o |
|        S+B   .  |
|        o+.*     |
|       o oo o    |
|        +  .     |
|         E.      |
+----[SHA256]-----+

$ ssh-copy-id -i id_rsa.pub root@www.cloudzp.club
````
输入密码即可
