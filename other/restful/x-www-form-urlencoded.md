## x-www-form-urlencoded

### 前言

最近使用`x-www-form-urlencoded`，同时也使用了加密的方式传递数据，`x-www-form-urlencoded`对数据进行了encoded，
在测试的时候使用encoded之后的数据就行测试，来来回回浪费了好多时间。  

### 原理探究

请求的时候，请求头中的Content-Type指明了主体部分的MIME类型，当类型是`application/x-www-form-urlencoded`  



