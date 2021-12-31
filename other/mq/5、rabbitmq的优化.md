## RabbitMQ 的优化

### channel 理解

生产者，消费者和 RabbitMQ 都会建立连接。为了避免建立过多的 TCP 连接，减少资源额消耗。  

AMQP 协议引入了信道(channel)，多个 channel 使用同一个 TCP 连接，起到对 TCP 连接的复用。    

不过 channel 的连接数是有上限的，过多的连接会导致复用的 TCP 拥堵。   