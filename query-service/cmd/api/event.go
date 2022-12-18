package main

import (
	"log"

	"github.com/Adhiana46/query-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	articleCreatedEvent = "article.created"
	articleUpdatedEvent = "article.updated"
	articleDeletedEvent = "article.deleted"
)

func (app *Config) listenEvents(topic string, events []string) {
	// create consumer
	consumer, err := event.NewConsumer(app.rabbitConn, topic, app.handleEvent)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen(events)
	if err != nil {
		log.Println(err)
	}
}

func (app *Config) handleEvent(msg *amqp.Delivery) {
	switch msg.RoutingKey {
	case articleCreatedEvent:
		log.Println("Article Created")
	case articleUpdatedEvent:
		log.Println("Article Updated")
	case articleDeletedEvent:
		log.Println("Article Deleted")
	}
}
