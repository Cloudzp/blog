# MY BLOG
## 启动

### 1. 执行打包任务
```
$ gulp 或者 gulp default
```
### 2. 清空hexo静态文件和缓存，并重新生成
```
$ hexo clean && hexo g
```
### 3. 本地预览，确没有问题再进行发布
```
$ hexo s -p 4000 或者 hexo s
```
### 4. 当gulp执行完成，并提示 please execute： hexo d 时，可以进行发布
```
$ hexo d 或者 gulp deploy  //部署发布
```

## 写文章
```
# hexo new "【文章主题】"
```
然后在生成的文件下进行编写，markdown格式即可: