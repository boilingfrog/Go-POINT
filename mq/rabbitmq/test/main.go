package main

import (
	"Go-POINT/mq/rabbitmq"
	"Go-POINT/mq/rabbitmq/help"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DeadTestExchangeQueue = "dead-test-delayed-queue_queue"
)

func main() {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	broker := rabbitmq.NewBroker("amqp://admin:admin@127.0.0.1:5672", &rabbitmq.ExchangeConfig{
		Name: "worker-exchange",
		Type: "direct",
	})

	broker.LaunchJobs(
		rabbitmq.NewDefaultJobber(
			"dead-test-key",
			HandleMessage,
			rabbitmq.WithPrefetch(30),
			rabbitmq.WithQueue(DeadTestExchangeQueue),
			rabbitmq.WithRetry(help.FIBONACCI, help.Retry{
				Delay: "5s",
				Max:   6,
				Queue: []string{
					DeadTestExchangeQueue,
				},
			}),
		),
	)

	for {
		s := <-ch
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			fmt.Println("job-test-exchange service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func HandleMessage(data []byte) error {
	fmt.Println("receive message", "message", string(data))

	return rabbitmq.HandleRequeue
}
