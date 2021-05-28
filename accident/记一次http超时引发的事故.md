<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [记一次http超时引发的事故](#%E8%AE%B0%E4%B8%80%E6%AC%A1http%E8%B6%85%E6%97%B6%E5%BC%95%E5%8F%91%E7%9A%84%E4%BA%8B%E6%95%85)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [分下下http client超时的处理方法](#%E5%88%86%E4%B8%8B%E4%B8%8Bhttp-client%E8%B6%85%E6%97%B6%E7%9A%84%E5%A4%84%E7%90%86%E6%96%B9%E6%B3%95)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## 记一次http超时引发的事故

### 前言

我们使用的是golang标准库的`http client`，对于一些http请求，我们在处理的时候，会考虑加上超时时间，防止http请求一直在请求，导致业务长时间阻塞等待。  

### 分下下http client超时的处理方法






### 参考

【[译]Go net/http 超时机制完全手册】https://colobu.com/2016/07/01/the-complete-guide-to-golang-net-http-timeouts/  