package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel, name string) error {
	return ch.ExchangeDeclare(
		name,    // name exchange
		"topic", // type
		true,    // durable?
		false,   // auto-delete?
		false,   // use-internally?
		false,   // no-wait?
		nil,     // arguments
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
