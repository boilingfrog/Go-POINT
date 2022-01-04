<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [RabbitMQ 的优化](#rabbitmq-%E7%9A%84%E4%BC%98%E5%8C%96)
  - [channel](#channel)
  - [prefetch Count](#prefetch-count)
  - [死信队列](#%E6%AD%BB%E4%BF%A1%E9%98%9F%E5%88%97)
    - [什么是死信队列](#%E4%BB%80%E4%B9%88%E6%98%AF%E6%AD%BB%E4%BF%A1%E9%98%9F%E5%88%97)
    - [使用场景](#%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF)
    - [代码实现](#%E4%BB%A3%E7%A0%81%E5%AE%9E%E7%8E%B0)
  - [延迟队列](#%E5%BB%B6%E8%BF%9F%E9%98%9F%E5%88%97)
    - [什么是延迟队列](#%E4%BB%80%E4%B9%88%E6%98%AF%E5%BB%B6%E8%BF%9F%E9%98%9F%E5%88%97)
    - [使用场景](#%E4%BD%BF%E7%94%A8%E5%9C%BA%E6%99%AF-1)
    - [代码实现](#%E4%BB%A3%E7%A0%81%E5%AE%9E%E7%8E%B0-1)
  - [参考](#%E5%8F%82%E8%80%83)

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

所以消息不会被处理速度很慢的消费者过多霸占，能够很好的分配到其它处理速度较好的消费者中。通俗的说就是消费者最多从 RabbitMQ 中获取的未消费消息的数量。          

`prefetch Count`数量设置为多少合适呢？大概就是30吧，具体可以参见[Finding bottlenecks with RabbitMQ 3.3](https://blog.rabbitmq.com/posts/2014/04/finding-bottlenecks-with-rabbitmq-3-3)  

谈到了`prefetch Count`，我们还要看了 global 这个参数,RabbitMQ 为了提升相关的性能，在` AMQPO-9-1` 协议之上重新定义了 global 这个参数  

| global 参数 |         AMQPO-9-1                                               | RabbitMQ |
| ------     | ------------------------------------------                       | ------------------------------------------------ |
| false      | 信道上所有的消费者都需要遵从 prefetchC unt 的限                       | 信道上新的消费者需要遵从 prefetchCount 的限定值定值 |
| true       | 当前通信链路( Connection) 上所有的消费者都要遵从 prefetchCount 的限定值 | 信道上所有的消费者都需要遵从 prefetchCunt 的上限，就是同一信道上的消费者共享 |

prefetchSize：预读取的单条消息内容大小上限(包含)，可以简单理解为消息有效载荷字节数组的最大长度限制，0表示无上限，单位为 B。   

如果`prefetch Count`为 0 呢，表示预读取的消息数量没有上限。

举个错误使用的栗子：  

之前一个队列的消费者消费速度过慢，`prefetch Count`为0，然后新写了一个消费者，`prefetch Count`设置为30，并且起了10个pod，来处理消息。老的消费者还没有下线也在处理消息。     

但是发现消费速度还是很慢，有大量的消息处于 unacked 。如果明白`prefetch Count`的含义其实就已经可以猜到问题的原因了。   

老的消费者`prefetch Count`为0，所以很多 unacked 消息都被它持有了，虽然新加了几个新的消费者，但是都处于空闲状态，最后停掉了`prefetch Count`为0的消费者，很快消费速度就正常了。   

### 死信队列

#### 什么是死信队列

一般消息满足下面几种情况就会消息变成死信  

- 消息被否定确认，使用 `channel.basicNack` 或 `channel.basicReject` ，并且此时 requeue 属性被设置为false； 

- 消息过期，消息在队列的存活时间超过设置的 TT L时间；  

- 队列达到最大长度，消息队列的消息数量已经超过最大队列长度。   

当一个消息满足上面的几种条件变成死信(dead message)之后，会被重新推送到死信交换器(DLX ，全称为 Dead-Letter-Exchange)。绑定 DLX 的队列就是私信队列。   

所以死信队列也并不是什么特殊的队列，只是绑定到了死信交换机中了，死信交换机也没有什么特殊，我们只是用这个来处理死信队列了，和别的交换机没有本质上的区别。   

对于需要处理私信队列的业务，跟我们正常的业务处理一样，也是定义一个独有的路由key，并对应的配置一个死信队列进行监听，然后 key 绑定的死信交换机中。   

#### 使用场景

当消息的消费出现问题时，出问题的消息不被丢失，进行消息的暂存，方便后续的排查处理。    

#### 代码实现


### 延迟队列

#### 什么是延迟队列

延迟队列就是用来存储进行延迟消费的消息。  

什么是延迟消息？   

就是不希望消费者马上消费的消息，等待指定的时间才进行消费的消息。     

#### 使用场景

1、关闭空闲连接。服务器中，有很多客户端的连接，空闲一段时间之后需要关闭之；  

2、清理过期数据业务上。比如缓存中的对象，超过了空闲时间，需要从缓存中移出；  

3、任务超时处理。在网络协议滑动窗口请求应答式交互时，处理超时未响应的请求；  

4、下单之后如果三十分钟之内没有付款就自动取消订单；  

5、订餐通知:下单成功后60s之后给用户发送短信通知；  

6、当订单一直处于未支付状态时，如何及时的关闭订单，并退还库存；  

7、如何定期检查处于退款状态的订单是否已经退款成功；  

8、新创建店铺，N天内没有上传商品，系统如何知道该信息，并发送激活短信；  

9、定时任务调度：使用DelayQueue保存当天将会执行的任务和执行时间，一旦从DelayQueue中获取到任务就开始执行。   

总结下来就是一些延迟处理的业务场景  

#### 代码实现

RabbitMQ 中本身并没有直接延迟队列的功能，可以通过死信队列和 TTL 。来实现延迟队的功能。   




### 参考

【Finding bottlenecks with RabbitMQ 3.3】https://blog.rabbitmq.com/posts/2014/04/finding-bottlenecks-with-rabbitmq-3-3  
【你真的了解延时队列吗】https://juejin.cn/post/6844903648397525006    


