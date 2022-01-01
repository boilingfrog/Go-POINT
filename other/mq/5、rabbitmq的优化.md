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

什么是`prefetch Count`，先举个栗子：  

假定 RabbitMQ 队列有 N 个消费队列，RabbitMQ 队列中的消息将以轮询的方式发送给消费者。   

消息的数量是 M,那么每个消费者得到的数据就是 M%N。如果某一台的机器中的消费者，因为自身的原因，或者消息本身处理所需要的时间很久，消费的很慢，但是其他消费者分配的消息很快就消费完了，然后处于闲置状态，这就造成资源的浪费，消息队列的吞吐量也降低了。   

这时候`prefetch Count`就登场了，通过引入`prefetch Count`来避免消费能力有限的消息队列分配过多的消息，而消息处理能力较好的消费者没有消息处理的情况。   

RabbitM 会保存一个消费者的列表，每发送一条消息都会为对应的消费者计数，如果达到了所设定的上限，那么 RabbitMQ 就不会向这个消费者再发送任何消息。直到消费者确认了某条消息之后 RabbitMQ 将相应的计数减1，之后消费者可以继续接收消息，直到再次到达计数上限。这种机制可以类比于 TCP!IP中的"滑动窗口"。   

通俗的说就是消费者最多从 RabbitMQ 中获取的未消费消息的数量。     

`prefetch Count`数量设置为多少合适呢？大概就是30吧，具体可以参见[Finding bottlenecks with RabbitMQ 3.3](https://blog.rabbitmq.com/posts/2014/04/finding-bottlenecks-with-rabbitmq-3-3)  

### 参考

【Finding bottlenecks with RabbitMQ 3.3】https://blog.rabbitmq.com/posts/2014/04/finding-bottlenecks-with-rabbitmq-3-3  


