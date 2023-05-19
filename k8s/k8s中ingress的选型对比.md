<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [k8s 中的选型对比](#k8s-%E4%B8%AD%E7%9A%84%E9%80%89%E5%9E%8B%E5%AF%B9%E6%AF%94)
  - [前言](#%E5%89%8D%E8%A8%80)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## k8s 中的选型对比

### 前言

最近涉及到业务迁移，梳理了老业务，发现用到几种不同的 ingress，借着这个机会来对比下市面中几种常用 ingress 的优缺点，以及业务中该如何选型。   

首先这里列举下，常用的几种 ingress。   

1、`Ingress Nginx`；  

2、`Nginx Ingress`；  

3、`Kong Ingress`；  

4、`Traefik Ingress`；  

5、`HAProxy Ingress`；  

6、`Istio Ingress`；  

7、`APISIX Ingress`；  

### Ingress Nginx

项目地址[https://github.com/kubernetes/ingress-nginx]   

`Ingress Nginx` 是官方推荐的 Ingress 控制器。它基于 Nginx 的平台，，并补充了一组用于实现额外功能的 Lua 插件。   

优点

1、因为是基于 Nginx 平台，Nginx 现在是世界上最流行的 `Nginx HTTP Sever`，由于 Nginx 的普及，该控制器最容易上手；   

2、同时接入 K8s 中的配置项也是最少的，同时文档资料也比较全，学习成本低，对于部分刚接触 K8S 的人或者创业公司来说，`Ingress Nginx` 确实是一个非常好的选择。  

缺点  

1、`Nginx reload` 耗时问题，`Ingress Nginx` 使用到了一些 OpenResty 的特性，但最终配置加载还是依赖于原有的 `Nginx config reload`。当路由配置非常大的时候，`Nginx reload` 会耗时非常久，可以达到几秒甚至十几秒，这种 reload 会很严重的影响业务，甚至造成业务中断；   

2、插件的开发困难，`Ingress Nginx` 本身提供的插件如果不能满足需求，定制化的开发会比较麻烦。  

### Nginx Ingress

项目地址[https://github.com/nginxinc/kubernetes-ingress] 

`Nginx Ingress` 是有 Nginx 官方维护的项目。  

优点  

1、NGINX 控制器具有长期稳定性和与一致性的特点，有很强的持续兼容性；  

2、剔除了第三方模块和 Lua 脚本；  

缺点  

1、

### 参考

【NGINX Ingress Controller 选项】https://www.nginx-cn.net/blog/guide-to-choosing-ingress-controller-part-4-nginx-ingress-controller-options/