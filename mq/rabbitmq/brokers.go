package rabbitmq

import (
	"Go-POINT/mq/rabbitmq/help"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func NewDefaultJober(key string, functor func([]byte) error, params ...Param) *params {
	var defaultParam = []Param{
		WithConcurrency(5),
		WithPrefetch(10),
		WithRetry(help.FIBONACCI, help.Retry{
			Delay: "2s",
			Max:   5,
			Queue: nil,
		}),
	}
	ps := evaParam(append(defaultParam, params...))
	ps.key = key
	ps.handler = DefaultHandler(key, functor)
	return ps
}

func (b *Broker) LaunchJobs(jobs ...*params) {
	var num = len(jobs)
	if num < 1 {
		return
	}

	for i := range jobs {
		go func(i int) {
			job := jobs[i]
			b.launchJobs(job)
		}(i)
	}
}

func (b *Broker) launchJobs(ps *params) {
	var key = ps.key

	for {
		log.Printf("job %s starting...", key)
		retry, err := b.readyConsumes(ps)
		if retry {
			log.Printf("job %s failed with error: %s, retrying start after 30 seconds...", key, err)
		}
	}
}

func (b *Broker) readyConsumes(ps *params) (bool, error) {
	var (
		key = ps.key
	)

	channel, err := b.getChannel(key)
	if err != nil {
		return true, err
	}

	queue, err := b.declare(channel, key, ps)
	if err != nil {
		return true, err
	}

	if err := channel.Qos(ps.prefetch, 0, false); err != nil {
		return true, fmt.Errorf("channel qos error: %s", err)
	}

	deliveries, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return true, fmt.Errorf("queue consume error: %s", err)
	}

	channelClose := channel.NotifyClose(make(chan *amqp.Error))

	pool := make(chan struct{}, ps.concurrency)

	go func() {
		for i := 0; i < ps.concurrency; i++ {
			pool <- struct{}{}
		}
	}()

	for {
		select {
		case err := <-channelClose:
			b.channels.Delete(key)
			return true, fmt.Errorf("channel close: %s", err)
		case d := <-deliveries:
			if ps.concurrency > 0 {
				<-pool
			}
			go func() {
				var flag HandleFLag

				switch flag = ps.Handle(d.Body); flag {
				case HandleSuccess:
					d.Ack(false)
				case HandleDrop:
					d.Nack(false, false)
				case HandleRequeue:
					if err := b.retry(ps, d); err != nil {
						d.Nack(false, true)
					} else {
						d.Ack(false)
					}
				default:
					d.Nack(false, false)
				}

				if ps.concurrency > 0 {
					pool <- struct{}{}
				}
			}()
		}
	}
}

func (b *Broker) retry(ps *params, d amqp.Delivery) error {
	channel, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	retryCount, _ := d.Headers["x-retry-count"].(int32)
	if int(retryCount) >= len(ps.retryQueue) {
		return nil
	}

	delay := ps.retryQueue[retryCount]
	delayDuration := time.Duration(delay) * time.Millisecond
	delayQ := fmt.Sprintf("delay.%s.%s.%s", delayDuration.String(), b.exchange, ps.key)

	if _, err := channel.QueueDeclare(delayQ,
		true, false, false, false, amqp.Table{
			"x-dead-letter-exchange":    b.exchange,
			"x-dead-letter-routing-key": ps.key,
			"x-message-ttl":             delay,
			"x-expires":                 delay * 2,
		},
	); err != nil {
		return err
	}

	return channel.Publish("", delayQ, false, false, amqp.Publishing{
		Headers:      amqp.Table{"x-retry-count": retryCount + 1},
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
	})
}
