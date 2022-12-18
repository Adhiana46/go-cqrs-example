package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	appName    = "Query Service"
	appVersion = "1.0"
	port       = "80"
)

type Config struct {
	AppName    string
	AppVersion string

	rabbitConn *amqp.Connection
}

func main() {
	app := Config{
		AppName:    appName,
		AppVersion: appVersion,
	}

	// open rabbitmq
	err := app.openRabbitmq()
	if err != nil {
		log.Panicf("Can't open RabbitMQ connection: %s", err)
	}
	defer app.closeRabbitmq()

	log.Printf("Starting %s service on port %s\n", appName, port)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// listening for events
	go app.listenEvents("articles", []string{"article.created", "article.updated", "article.deleted"})

	// starting the server
	if err := s.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

// Rabbitmq
func (app *Config) openRabbitmq() error {
	var count int64
	var retryTime = 1 * time.Second
	var connection *amqp.Connection

	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("AMQP_USER"),
		os.Getenv("AMQP_PASSWORD"),
		os.Getenv("AMQP_HOST"),
		os.Getenv("AMQP_PORT"),
	)

	// Don't continue until rabbit is ready
	for {
		c, err := amqp.Dial(dsn)
		if err != nil {
			log.Println("RabbitMQ not yet ready...", err)
			count++
		} else {
			log.Println("Connected to RabbitMQ...")
			connection = c
			break
		}

		if count > 5 {
			log.Println("Could not connect to RabbitMQ", err)
			return err
		}

		retryTime = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("Retrying in", retryTime)
		time.Sleep(retryTime)
		continue
	}

	app.rabbitConn = connection

	return nil
}

func (app *Config) closeRabbitmq() {
	app.rabbitConn.Close()
}
