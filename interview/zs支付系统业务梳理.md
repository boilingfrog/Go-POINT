## 聊一聊支付系统如何做到不丢单

### 前言

这里来聊一聊公司在支付系统中的技术选型。我们充值业务相对来说较为简单，充值的商品只是 VIP , 书币 这种相对单一的商品，没有涉及到电商业务中的购物车，合并支付这种复杂的支付场景。  

所以总结下来，主要有一下几个点：    

1、安全，避免被刷；  

2、不丢单，在异常的网络情况下，能购保证充值正常到账；  

3、并发，如果充值瞬时流量很大，业务系统是否能抗的住；  

4、简单易扩展，因为目前的支付有微信，支付宝，银联。。。各种支付，如果要接入这几种，支付系统的设计过程中，要易于扩展。  

好了，下面介绍下具体的实现  

### 分布式系统

不管是支付系统还是其他的系统，目前的设计已经摒弃了单机，采用的都是分布式的部署方案。  

使用分布式系统有点很明显，扩展性，并发性能顾



