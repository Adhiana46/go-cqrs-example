package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	exchangeName string
	connection   *amqp.Connection
}

func (e *Emitter) setup() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch, e.exchangeName)
}

func (e *Emitter) Push(eventName string, data []byte) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// declare queue
	queue, err := declareQueue(ch, e.exchangeName)
	if err != nil {
		return err
	}

	log.Printf("Push to channel E:%s -> R:%s -> Q:%s", e.exchangeName, eventName, queue.Name)

	ch.QueueBind(queue.Name, eventName, e.exchangeName, false, nil)

	err = ch.Publish(
		e.exchangeName,
		eventName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventEmitter(conn *amqp.Connection, exchangeName string) (Emitter, error) {
	emitter := Emitter{
		connection:   conn,
		exchangeName: exchangeName,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
