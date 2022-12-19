package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/Adhiana46/query-service/query"
	"github.com/go-redis/redis/v9"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	appName    = "Query Service"
	appVersion = "1.0"
	port       = "80"
)

type Config struct {
	AppName    string
	AppVersion string

	mongoDb    *mongo.Client
	rabbitConn *amqp.Connection
	rds        *redis.Client

	queryArticle query.ArticleQuery
}

func main() {
	app := Config{
		AppName:    appName,
		AppVersion: appVersion,
	}

	// open mongodb
	err := app.openMongodb()
	if err != nil {
		log.Panicf("Can't open MongoDB connection: %s", err)
	}
	defer app.closeMongodb()

	// open rabbitmq
	err = app.openRabbitmq()
	if err != nil {
		log.Panicf("Can't open RabbitMQ connection: %s", err)
	}
	defer app.closeRabbitmq()

	// open redis
	err = app.openRedis()
	if err != nil {
		log.Panicf("Can't open Redis connection: %s", err)
	}
	defer app.closeRedis()

	app.registerQuery()

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

func (app *Config) registerQuery() {
	app.queryArticle = query.NewArticleQueryMongo(app.mongoDb, app.rabbitConn, app.rds)
}

// Mongodb
func (app *Config) openMongodb() error {
	mongoURL := os.Getenv("MONGO_URL")
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")

	var count int64
	var retryTime = 1 * time.Second

	for {
		clientOptions := options.Client().ApplyURI(mongoURL)
		clientOptions.SetAuth(options.Credential{
			Username: username,
			Password: password,
		})

		c, err := mongo.Connect(context.TODO(), clientOptions)

		if err != nil {
			log.Println("MongoDB not yet ready...", err)
			count++
		} else {
			log.Println("Connected to MongoDB...")
			app.mongoDb = c
			break
		}

		if count > 5 {
			log.Println("Could not connect to MongoDB", err)
			return err
		}

		retryTime = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("Retrying in", retryTime)
		time.Sleep(retryTime)
		continue
	}

	return nil
}

func (app *Config) closeMongodb() {
	//
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

// Redis
func (app *Config) openRedis() error {
	app.rds = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%v:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password:    os.Getenv("REDIS_PASSWORD"),
		DB:          0, // use default DB
		ReadTimeout: -1,
	})

	log.Println("Connected to Redis")

	return nil
}

func (app *Config) closeRedis() {
	app.rds.Close()
}
