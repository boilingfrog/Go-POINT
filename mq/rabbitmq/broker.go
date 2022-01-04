package rabbitmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/streadway/amqp"
)

type Broker struct {
	exchange     string
	exchangeType string
	url          string
	conn         *amqp.Connection
	connClose    chan *amqp.Error
	channels     sync.Map
	duration     *prometheus.SummaryVec
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

func (b *Broker) Monitoring(address string) {

	b.duration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "rabbitmq_job_summary",
		Help:       "rabbitmq job summary",
		Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.99: 0.001, 0.999: 0.0001},
	}, []string{"taskName", "state"})

	prometheus.MustRegister(b.duration)

	debugMux := http.NewServeMux()
	debugMux.Handle("/metrics", promhttp.Handler())
	debugMux.Handle("/healthz", b.healthCheck())

	go func() {
		log.Printf("monitoring server listen on port %s...\n", address)
		http.ListenAndServe(address, debugMux)
	}()
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

type Jober interface {
	Queue() string
	Handle([]byte) HandleFLag
}

func (b *Broker) LaunchJob(key string, concurrency int, job Jober) {
	for {
		log.Printf("job %s starting...", job.Queue())
		retry, err := b.readyConsume(key, concurrency, job)
		if retry {
			log.Printf("job %s failed with error: %s, retrying start after 30 seconds...", key, err)
			time.Sleep(30 * time.Second)
		}
	}
}

func (b *Broker) readyConsume(key string, concurrency int, job Jober) (bool, error) {
	channel, err := b.getChannel(key)
	if err != nil {
		return true, err
	}

	queue, err := b.declare(channel, key, job)
	if err != nil {
		return true, err
	}

	if err := channel.Qos(10, 0, false); err != nil {
		return true, fmt.Errorf("channel qos error: %s", err)
	}

	deliveries, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return true, fmt.Errorf("queue consume error: %s", err)
	}

	channelClose := channel.NotifyClose(make(chan *amqp.Error))

	pool := make(chan struct{}, concurrency)

	go func() {
		for i := 0; i < concurrency; i++ {
			pool <- struct{}{}
		}
	}()

	for {
		select {
		case err := <-channelClose:
			b.channels.Delete(key)
			return true, fmt.Errorf("channel close: %s", err)
		case d := <-deliveries:
			if concurrency > 0 {
				<-pool
			}
			go func() {
				var flag HandleFLag

				defer func(begin time.Time) {
					b.duration.With(prometheus.Labels{"taskName": key, "state": string(flag)}).
						Observe(float64(time.Since(begin) / time.Millisecond))
				}(time.Now())

				switch flag = job.Handle(d.Body); flag {
				case HandleSuccess:
					d.Ack(false)
				case HandleDrop:
					d.Nack(false, false)
				case HandleRequeue:
					d.Nack(false, true)
				default:
					d.Nack(false, false)
				}

				if concurrency > 0 {
					pool <- struct{}{}
				}
			}()
		}
	}
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

func (b *Broker) BatchPublish(key string, dataList []interface{}) error {
	channel, err := b.createChannel()
	if err != nil {
		return err
	}
	defer channel.Close()

	var body []byte

	for _, data := range dataList {
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

		if err := channel.Publish(b.exchange, key, false, false, amqp.Publishing{
			Headers:      amqp.Table{},
			ContentType:  "",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		}); err != nil {
			return err
		}
	}

	return nil
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

func (b *Broker) declare(channel *amqp.Channel, key string, job Jober) (amqp.Queue, error) {
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
