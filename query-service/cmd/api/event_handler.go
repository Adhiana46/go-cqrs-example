package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Adhiana46/query-service/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
)

func (app *Config) handleArticleCreated(msg *amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	article := model.Article{}
	_ = json.Unmarshal(msg.Body, &article)

	if article.Uuid != "" {
		// Set Cache
		cacheKey := fmt.Sprintf("article-%s", article.Uuid)
		app.rds.Set(ctx, cacheKey, string(msg.Body), 10*time.Minute)

		// insert into collection
		collection := app.mongoDb.Database("articles").Collection("articles")

		article.ID = ""
		_, _ = collection.InsertOne(ctx, article)
	}

	msg.Ack(false)
}

func (app *Config) handleArticleUpdated(msg *amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	article := model.Article{}
	_ = json.Unmarshal(msg.Body, &article)

	if article.Uuid != "" {
		// Set Cache
		cacheKey := fmt.Sprintf("article-%s", article.Uuid)
		app.rds.Set(ctx, cacheKey, string(msg.Body), 10*time.Minute)

		// update into collection
		collection := app.mongoDb.Database("articles").Collection("articles")
		_, _ = collection.UpdateOne(
			ctx,
			bson.M{"uuid": article.Uuid},
			bson.D{
				{Key: "$set", Value: bson.D{
					{Key: "uuid", Value: article.Uuid},
					{Key: "author", Value: article.Author},
					{Key: "title", Value: article.Title},
					{Key: "body", Value: article.Body},
					{Key: "created_at", Value: article.CreatedAt},
					{Key: "updated_at", Value: article.UpdatedAt},
				}},
			},
		)
	}

	msg.Ack(false)
}

func (app *Config) handleArticleDeleted(msg *amqp.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	article := model.Article{}
	_ = json.Unmarshal(msg.Body, &article)

	if article.Uuid != "" {
		// Delete Cache
		cacheKey := fmt.Sprintf("article-%s", article.Uuid)
		app.rds.Del(ctx, cacheKey)

		// Delete from collection
		collection := app.mongoDb.Database("articles").Collection("articles")

		_, _ = collection.DeleteOne(ctx, bson.M{"uuid": article.Uuid})
	}

	msg.Ack(false)
}
