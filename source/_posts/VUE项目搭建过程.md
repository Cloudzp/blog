---
title: VUE项目搭建过程
categories:
  - 前段
tags:
  - js
  - vue
comments: true
abbrlink: 64724
date: 2018-07-18 21:11:18
img:
---

> https://www.jianshu.com/p/02b12c600c7b
// TODO
## 安装npm
```
```
## VUE中的html格式设置？
[点击查看](https://vuetifyjs.com/zh-Hans/components/grids)

## vue代码中如何跳转到其他页面？
- html中跳转
````vue

<template>
   <v-btn :to="{path:'/app/create',query:{key:item}}">
   </v-btn>
</template>

````
- js中跳转
```javascript
 window.location.href = `/app/create`

```

## vue 如何校验form表单数据有效性？
````vue
<template>
   <v-from ref="form">
      <v-text-field 
      v-model="name"
      :rules="rules.required"
      >
      
      </v-text-field>
   </v-from>
</template>

<script>
  export default {
      name: "test",
      data(){
          return {
              rules: {                 
                  required: val => !!val || 'name为必填字段',
                  nameRule: val => {
                      const pat = /^[0-9a-z\-]{1,253}/;
                      return pat.test(val)
                  }
              }
          }
      }
  }
</script>

````



