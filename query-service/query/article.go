package query

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Adhiana46/query-service/dto"
	"github.com/Adhiana46/query-service/model"
	"github.com/go-redis/redis/v9"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleQuery interface {
	GetSingle(ctx context.Context, reqDto dto.RequestSingleArticle) (*model.Article, error)
	GetList(ctx context.Context, reqDto dto.RequestListArticle) ([]*model.Article, error)
}

type articleQueryMongo struct {
	mongoDb    *mongo.Client
	rabbitConn *amqp.Connection
	rds        *redis.Client
}

func NewArticleQueryMongo(mongoDb *mongo.Client, rabbitConn *amqp.Connection, rds *redis.Client) ArticleQuery {
	return &articleQueryMongo{
		mongoDb:    mongoDb,
		rabbitConn: rabbitConn,
		rds:        rds,
	}
}

func (query *articleQueryMongo) GetSingle(ctx context.Context, reqDto dto.RequestSingleArticle) (*model.Article, error) {
	cacheKey := fmt.Sprintf("article-%s", reqDto.Uuid)
	cacheResult, err := query.rds.Get(ctx, cacheKey).Result()
	var article model.Article

	// get from cache, ignore error
	if err == nil && cacheResult != "" {
		err = json.Unmarshal([]byte(cacheResult), &article)
		if err == nil {
			return &article, nil
		}
	}

	// get from mongodb collection
	collection := query.mongoDb.Database("articles").Collection("articles")

	err = collection.FindOne(ctx, bson.M{"uuid": reqDto.Uuid}).Decode(&article)
	if err != nil {
		return nil, err
	}

	// Store to cache
	articleJson, err := json.Marshal(article)
	if err != nil {
		return nil, err
	}

	err = query.rds.Set(ctx, cacheKey, string(articleJson), 10*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (query *articleQueryMongo) GetList(ctx context.Context, reqDto dto.RequestListArticle) ([]*model.Article, error) {
	// cache-key based on reqDto json -> md5
	reqDtoJson, _ := json.Marshal(reqDto)
	hash := md5.Sum(reqDtoJson)
	cacheKey := fmt.Sprintf("article-list-%s", hex.EncodeToString(hash[:]))

	var articles []*model.Article

	// get from cache, ignore error
	cacheResult, err := query.rds.Get(ctx, cacheKey).Result()
	if err == nil && cacheResult != "" {
		err = json.Unmarshal([]byte(cacheResult), &articles)
		if err == nil {
			return articles, nil
		}
	}

	collection := query.mongoDb.Database("articles").Collection("articles")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("Get list of articles error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var article model.Article

		err := cursor.Decode(&article)
		if err != nil {
			log.Println("Error decoding article into model.Article:", err)
		} else {
			articles = append(articles, &article)
		}
	}

	// Store to cache
	articlesJson, err := json.Marshal(articles)
	if err != nil {
		return nil, err
	}

	err = query.rds.Set(ctx, cacheKey, string(articlesJson), 10*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return articles, nil
}
