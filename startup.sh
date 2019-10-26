#!/bin/sh
gulp
hexo clean && hexo g
hexo s -p 4000 -ip 0.0.0.0