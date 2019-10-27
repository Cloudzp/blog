#!/bin/sh

## deploy to pro
### remove old source
ssh   root@www.cloudzp.club rm -rf /root/workdir/blog/source
scp -C -r source   root@www.cloudzp.club:/root/workdir/blog

### remove old _config
ssh   root@www.cloudzp.club rm -rf /root/workdir/blog/themes/hexo-theme-snippet/_config.yml
scp -C themes/hexo-theme-snippet/_config.yml root@www.cloudzp.club:/root/workdir/blog/themes/hexo-theme-snippet

scp -C -r themes/hexo-theme-snippet/source root@www.cloudzp.club:/root/workdir/blog/themes/hexo-theme-snippet

### restart server
ssh  root@www.cloudzp.club docker restart 94c9433f58ec

if [ "$?" != "0" ]; then 
   echo "deploy failed !"
   exit 1
fi
echo "deploy success !"

## git push
git add ./*
git commit -m "edit $(date)"
git push origin master