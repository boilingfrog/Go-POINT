package rabbitmq

import (
	"Go-POINT/mq/rabbitmq/help"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Broker struct {
	exchange     string
	exchangeType string
	url          string
	conn         *amqp.Connection
	connClose    chan *amqp.Error
	channels     sync.Map
}

type ExchangeConfig struct {
	Name string
	Type string
}

func NewBroker(url string, cfg *ExchangeConfig) *Broker {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	if cfg.Type == "" {
		cfg.Type = "direct"
	}

	var broker = &Broker{
		url:          url,
		conn:         conn,
		exchange:     cfg.Name,
		exchangeType: cfg.Type,
		connClose:    conn.NotifyClose(make(chan *amqp.Error)),
	}

	go func() {
		<-broker.connClose
		broker.conn = nil
	}()

	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			broker.checkConnection()
		}
	}()

	return broker
}

func (b *Broker) checkConnection() {
	if b.conn != nil {
		return
	}
	conn, err := amqp.Dial(b.url)

	if err != nil {
		log.Printf("broker redial faild: %v", err)
		return
	}
	b.connClose = conn.NotifyClose(make(chan *amqp.Error))
	go func() {
		<-b.connClose
		b.conn = nil
	}()

	b.conn = conn
}

func (b *Broker) healthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		status := "UP"

		if b.conn == nil {
			status = "DOWN"
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(status))
	})
}

type HandleFLag string

func (h HandleFLag) Error() string {
	return string(h)
}

const (
	HandleSuccess HandleFLag = "success"
	HandleDrop    HandleFLag = "drop"
	HandleRequeue HandleFLag = "requeue"
)

type Jobber interface {
	Queue() string
	Handle([]byte) HandleFLag
}

func NewDefaultJobber(key string, functor func([]byte) error, params ...Param) *params {
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
	key := ps.key
	channel, err := b.getChannel(key)
	if err != nil {
		return true, err
	}

	b.declareDelay(channel, key, ps)

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

	fmt.Println(retryCount)
	if int(retryCount) >= len(ps.retryQueue) {
		return nil
	}
	fmt.Println("xxxx", retryCount)

	delay := ps.retryQueue[retryCount]
	//delayQ := fmt.Sprintf("delay.%s.%s", b.exchange, ps.key)

	if err := channel.ExchangeDeclare(fmt.Sprintf("delay.%s", b.exchange), "x-delayed-message", true, false, false, false, amqp.Table{
		"x-delayed-type": "direct",
	}); err != nil {
		return err
	}

	//if _, err := channel.QueueDeclare(delayQ,
	//	true, false, false, false, amqp.Table{
	//		"x-dead-letter-exchange":    b.exchange,
	//		"x-dead-letter-routing-key": ps.key,
	//	},
	//); err != nil {
	//	return err
	//}
	//
	//if err = channel.QueueBind(delayQ, "test-delayed-message-rou", "test-delayed-message-1", false, nil); err != nil {
	//	return fmt.Errorf("queue bind error: %s", err)
	//}

	return channel.Publish(fmt.Sprintf("delay.%s", b.exchange), fmt.Sprintf("delay.%s", ps.key), false, false, amqp.Publishing{
		Headers: amqp.Table{"x-retry-count": retryCount + 1},
		//Headers:      amqp.Table{"x-retry-count": retryCount + 1, "x-delay": 1000},
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
		Expiration:   fmt.Sprintf("%d", delay),
	})
}

func (b *Broker) createChannel() (*amqp.Channel, error) {
	if b.conn == nil {
		return nil, errors.New("no available connection when create channel")
	}
	channel, err := b.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("create channel error: %s", err)
	}
	if err = channel.Confirm(false); err != nil {
		return nil, fmt.Errorf("channel not into confirm mode: %s", err)
	}
	return channel, nil
}

func (b *Broker) getChannel(key string) (*amqp.Channel, error) {
	if value, ok := b.channels.Load(key); ok {
		if c, ok := value.(*amqp.Channel); ok {
			return c, nil
		}
	}
	channel, err := b.createChannel()
	if err != nil {
		return nil, err
	}
	b.channels.Store(key, channel)
	return channel, nil
}

func (b *Broker) declareDelay(channel *amqp.Channel, key string, job Jobber) (amqp.Queue, error) {
	if err := channel.ExchangeDeclare(fmt.Sprintf("delay.%s", b.exchange), "x-delayed-message", true, false, false, false, nil); err != nil {
		return amqp.Queue{}, fmt.Errorf("exchange declare error: %s", err)
	}

	queue, err := channel.QueueDeclare(fmt.Sprintf("delay.%s", job.Queue()), true, false, false, false, nil)
	if err != nil {
		return queue, fmt.Errorf("queue declare error: %s", err)
	}
	if err = channel.QueueBind(queue.Name, fmt.Sprintf("delay.%s", key), fmt.Sprintf("delay.%s", b.exchange), false, nil); err != nil {
		return queue, fmt.Errorf("queue bind error: %s", err)
	}
	return queue, nil
}

func (b *Broker) declare(channel *amqp.Channel, key string, job Jobber) (amqp.Queue, error) {
	if err := channel.ExchangeDeclare(b.exchange, b.exchangeType, true, false, false, false, nil); err != nil {
		return amqp.Queue{}, fmt.Errorf("exchange declare error: %s", err)
	}

	queue, err := channel.QueueDeclare(job.Queue(), true, false, false, false, nil)
	if err != nil {
		return queue, fmt.Errorf("queue declare error: %s", err)
	}
	if err = channel.QueueBind(queue.Name, key, b.exchange, false, nil); err != nil {
		return queue, fmt.Errorf("queue bind error: %s", err)
	}
	return queue, nil
}

func (b *Broker) Publish(key string, data interface{}) error {

	var err error
	channel, err := b.createChannel()
	if err != nil {
		return err
	}
	defer channel.Close()

	var body []byte
	if d, ok := data.(string); ok {
		body = []byte(d)
	} else {
		body, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	if err := channel.ExchangeDeclare(b.exchange, b.exchangeType, true, false, false, false, nil); err != nil {
		return err
	}

	return channel.Publish(b.exchange, key, false, false, amqp.Publishing{
		Headers:      amqp.Table{},
		ContentType:  "",
		Body:         body,
		DeliveryMode: amqp.Persistent,
	})
}
