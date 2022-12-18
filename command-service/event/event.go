package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel, exchangeName string) error {
	return ch.ExchangeDeclare(
		exchangeName, // name exchange
		"topic",      // type
		true,         // durable?
		false,        // auto-delete?
		false,        // use-internally?
		false,        // no-wait?
		nil,          // arguments
	)
}

func declareQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,  // name?
		true,  // durable?
		false, // delete when unuse?
		false, // exclusive?
		false, // no-wait?
		nil,   // args
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unuse?
		true,  // exclusive?
		false, // no-wait?
		nil,   // args
	)
}
