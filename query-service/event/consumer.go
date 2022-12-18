package event

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn         *amqp.Connection
	exchangeName string

	handlePayload func(msg *amqp.Delivery)
}

func NewConsumer(conn *amqp.Connection, exchangeName string, handlePayload func(msg *amqp.Delivery)) (Consumer, error) {
	consumer := Consumer{
		conn:          conn,
		exchangeName:  exchangeName,
		handlePayload: handlePayload,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(ch, c.exchangeName)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareQueue(ch, c.exchangeName)
	if err != nil {
		return err
	}

	// set Qos
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		err = ch.QueueBind(
			q.Name,
			topic,
			c.exchangeName,
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name, // queue name
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for msg := range messages {
			log.Println("[MSG]:", msg.Exchange, msg.RoutingKey)
			c.handlePayload(&msg)
		}
	}()

	fmt.Printf("Waiting for messages [Exchange, Queue] [%s, %s]\n", c.exchangeName, q.Name)

	<-forever

	return nil
}
