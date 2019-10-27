#!/bin/sh
gulp
hexo clean && hexo g
hexo s -p 4000 -i 0.0.0.0
