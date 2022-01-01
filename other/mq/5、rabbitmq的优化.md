<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [RabbitMQ 的优化](#rabbitmq-%E7%9A%84%E4%BC%98%E5%8C%96)
  - [channel](#channel)
  - [prefetch Count](#prefetch-count)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## RabbitMQ 的优化

### channel 

生产者，消费者和 RabbitMQ 都会建立连接。为了避免建立过多的 TCP 连接，减少资源额消耗。  

AMQP 协议引入了信道(channel)，多个 channel 使用同一个 TCP 连接，起到对 TCP 连接的复用。    

不过 channel 的连接数是有上限的，过多的连接会导致复用的 TCP 拥堵。   

### prefetch Count  

`prefetch Count`：定义 Federation 内部缓存的消息条数，即在收到上游消息之后且在发送到下游之前缓存的消息条数。 

通俗的说就是消费者一次从 RabbitMQ 中获取的消息的数量。由于消费者自身处理消息的能力有限，当获取一定的消息之后，不希望队列中的消息再次推送过来，所以使用了`prefetch Count`，当拉取的数据处理完，之后才会再次从队列中拉取数据。     

`prefetch Count`数量设置为多少合适呢？大概就是30吧，具体可以参见[Finding bottlenecks with RabbitMQ 3.3](https://blog.rabbitmq.com/posts/2014/04/finding-bottlenecks-with-rabbitmq-3-3)  



### 参考

【Finding bottlenecks with RabbitMQ 3.3】https://blog.rabbitmq.com/posts/2014/04/finding-bottlenecks-with-rabbitmq-3-3  


