package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/Adhiana46/command-service/command"
	"github.com/go-redis/redis/v9"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	appName    = "Command Service"
	appVersion = "1.0"
	port       = "80"
)

type Config struct {
	AppName    string
	AppVersion string

	DB         *sqlx.DB
	rabbitConn *amqp.Connection
	rds        *redis.Client

	cmdArticle command.ArticleCommand
}

func main() {
	app := Config{
		AppName:    appName,
		AppVersion: appVersion,
	}

	// open db connection (postgresql)
	err := app.openDB()
	if err != nil {
		log.Panicf("Can't open database connection: %s", err)
	}
	defer app.closeDB()

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

	app.registerCommand()

	log.Printf("Starting %s service on port %s\n", appName, port)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	// starting the server
	if err := s.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func (app *Config) registerCommand() {
	app.cmdArticle = command.NewArticleCommandPg(app.DB, app.rabbitConn, app.rds)
}

// Postgresql
func (app *Config) openDB() error {
	var count int64
	var retryTime = 1 * time.Second
	var connection *sqlx.DB

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("CMD_DB_HOST"),
		os.Getenv("CMD_DB_PORT"),
		os.Getenv("CMD_DB_USER"),
		os.Getenv("CMD_DB_DATABASE"),
		os.Getenv("CMD_DB_PASSWORD"),
	)

	for {
		c, err := sqlx.Connect("pgx", dsn)
		if err != nil {
			log.Println("Postgresql not ready yet...", err)
			count++
		} else {
			log.Println("Connected to Postgresql")

			c.SetMaxOpenConns(60)
			c.SetConnMaxLifetime(120 * time.Second)
			c.SetMaxIdleConns(30)
			c.SetConnMaxIdleTime(20 * time.Second)
			if err = c.Ping(); err != nil {
				return err
			}

			connection = c

			break
		}

		if count > 5 {
			log.Println("Could not connect to Postgresql", err)
			return err
		}

		retryTime = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("Retrying in", retryTime)
		time.Sleep(retryTime)
		continue
	}

	app.DB = connection

	return nil
}

func (app *Config) closeDB() {
	app.DB.Close()
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

	return nil
}

func (app *Config) closeRedis() {
	app.rds.Close()
}
