package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ExitCh = make(chan struct{})

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

const verbose = false

var deliveryCount = 0

const queueName = "oengus.bot" // Also our routing key

var consumer *Consumer

func SetupAMQP() *Consumer {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     "setup",
		done:    make(chan error),
	}

	config := amqp.Config{
		Vhost:      "/",
		Properties: amqp.NewConnectionProperties(),
	}
	config.Properties.SetClientConnectionName("timers")

	uri := os.Getenv("RABBIT_MQ_URI")

	if uri == "" {
		c.tag = "no-amqp-uri"
		return c
	}

	var err error

	log.Printf("producer: dialing %s", uri)

	c.conn, err = amqp.DialConfig(uri, config)
	if err != nil {
		log.Fatalf("producer: error in dial: %s", err)
	}

	log.Println("producer: got Connection, getting Channel")
	c.channel, err = c.conn.Channel()

	if err != nil {
		log.Fatalf("error getting a channel: %s", err)
	}

	c.tag = "ready"
	consumer = c

	return c
}

func PublishBotMessage(jsonBody string) {
	if consumer == nil || consumer.tag != "ready" {
		// TODO: handle error state
		return
	}

	err := consumer.channel.Publish(
		"amq.topic",
		queueName,
		true,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "UTF-8",
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
			AppId:           "oenus-timers",
			Body:            []byte(jsonBody),
		},
	)

	if err != nil {
		log.Fatalf("producer: error in publish: %s", err)
	}
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func SetupCloseHandler(consumer *Consumer) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Ctrl+C pressed in Terminal")
		if err := consumer.Shutdown(); err != nil {
			log.Fatalf("error during shutdown: %s", err)
		}
		os.Exit(0)
	}()
}
