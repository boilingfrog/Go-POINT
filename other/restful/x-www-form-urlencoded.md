<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [x-www-form-urlencoded](#x-www-form-urlencoded)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [使用](#%E4%BD%BF%E7%94%A8)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## x-www-form-urlencoded

### 前言

最近使用`x-www-form-urlencoded`，同时也使用了加密的方式传递数据，`x-www-form-urlencoded`对数据进行了encoded，
在测试的时候使用encoded之后的数据就行测试，来来回回浪费了好多时间。  

### 使用

请求的时候，请求头中的Content-Type指明了主体部分的MIME类型，当类型是`application/x-www-form-urlencoded`  

请求传递参数
![post_form](/img/post_form_1.jpg?raw=true)

会发生encoded
![post_form](/img/post_form.jpg?raw=true)

服务端正常接收就好了，不用再转译了  


