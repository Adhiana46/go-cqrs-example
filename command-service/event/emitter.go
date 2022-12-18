package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	topicName  string
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch, e.topicName)
}

func (e *Emitter) Push(eventName string, data []byte) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("Push to channel", e.topicName, eventName)

	err = ch.Publish(
		e.topicName,
		eventName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection, topicName string) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
		topicName:  topicName,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
