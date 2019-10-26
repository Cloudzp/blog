FROM node:latest
WORKDIR /blog
RUN cnpm install --global gulp-cli \
&& cnpm install gulp --save \
&& cnpm install -g hexo-cli
CMD [ "./startup.sh" ]