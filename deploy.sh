#!/bin/sh
## git push
git add ./*
git commit -m "edit $(date)"
git push origin master
## deploy to pro 
scp -r source   root@www.cloudzp.club:/root/workdir/blog
scp themes/hexo-theme-snippet/_config.yaml root@www.cloudzp.club:/root/workdir/blog/themes/hexo-theme-snippet
scp -r themes/hexo-theme-snippet/source root@www.cloudzp.club:/root/workdir/blog/themes/hexo-theme-snippet

echo "deploy success !"