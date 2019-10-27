FROM node:latest
WORKDIR /blog
RUN npm install --global cnpm \
&& cnpm install --global gulp-cli \
&& cnpm install --global gulp --save \
&& cnpm install --global hexo-cli
CMD [ "./startup.sh" ]
