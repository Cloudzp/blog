#!/bin/sh
## git push
git add ./*
git commit -m "edit $(date)"
git push origin master
## deploy to pro
ssh   root@www.cloudzp.club rm -rf /root/workdir/blog/source
scp -C -r source   root@www.cloudzp.club:/root/workdir/blog
scp -C themes/hexo-theme-snippet/_config.yml root@www.cloudzp.club:/root/workdir/blog/themes/hexo-theme-snippet
scp -C -r themes/hexo-theme-snippet/source root@www.cloudzp.club:/root/workdir/blog/themes/hexo-theme-snippet

echo "deploy success !"