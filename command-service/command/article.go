package command

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Adhiana46/command-service/dto"
	"github.com/Adhiana46/command-service/event"
	"github.com/Adhiana46/command-service/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	articleCreatedEvent = "article.created"
	articleUpdatedEvent = "article.updated"
	articleDeletedEvent = "article.deleted"
)

type ArticleCommand interface {
	PushToQueue(ctx context.Context, eventName string, article *model.Article) error
	Store(ctx context.Context, reqDto dto.RequestStoreArticle) (*model.Article, error)
	Update(ctx context.Context, reqDto dto.RequestUpdateArticle) (*model.Article, error)
	Delete(ctx context.Context, reqDto dto.RequestDeleteArticle) (*model.Article, error)
}

type articleCommandPg struct {
	db         *sqlx.DB
	rabbitConn *amqp.Connection
	rds        *redis.Client
}

func NewArticleCommandPg(db *sqlx.DB, rabbitConn *amqp.Connection, rds *redis.Client) ArticleCommand {
	return &articleCommandPg{
		db:         db,
		rabbitConn: rabbitConn,
		rds:        rds,
	}
}

func (c *articleCommandPg) PushToQueue(ctx context.Context, eventName string, article *model.Article) error {
	emitter, err := event.NewEventEmitter(c.rabbitConn, "articles")
	if err != nil {
		return err
	}

	jsonPayload, err := json.Marshal(article)
	if err != nil {
		return err
	}

	// Redis
	cacheKey := fmt.Sprintf("article-%s", article.Uuid)
	switch eventName {
	case articleCreatedEvent:
		c.rds.Set(ctx, cacheKey, article, 10*time.Minute)
	case articleUpdatedEvent:
		c.rds.Set(ctx, cacheKey, article, 10*time.Minute)
	case articleDeletedEvent:
		c.rds.Del(ctx, cacheKey)
	}

	return emitter.Push(eventName, jsonPayload)
}

func (c *articleCommandPg) findByUuid(ctx context.Context, uuid string) (*model.Article, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Select("*").
		From("articles").
		Where(sq.Eq{"uuid": uuid}).
		ToSql()

	if err != nil {
		return nil, err
	}

	row := model.Article{}
	err = c.db.GetContext(ctx, &row, sql, args...)
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (c *articleCommandPg) Store(ctx context.Context, reqDto dto.RequestStoreArticle) (*model.Article, error) {
	validate := validator.New()

	if err := validate.Struct(reqDto); err != nil {
		return nil, err
	}

	values := map[string]interface{}{
		"uuid":       uuid.NewString(),
		"author":     reqDto.Author,
		"title":      reqDto.Title,
		"body":       reqDto.Body,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	// Build Sql
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Insert("articles").
		SetMap(values).
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sql, args...)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	article, err := c.findByUuid(ctx, values["uuid"].(string))
	if err != nil {
		return nil, err
	}

	// Push
	err = c.PushToQueue(ctx, articleCreatedEvent, article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (c *articleCommandPg) Update(ctx context.Context, reqDto dto.RequestUpdateArticle) (*model.Article, error) {
	validate := validator.New()

	if err := validate.Struct(reqDto); err != nil {
		return nil, err
	}

	values := map[string]interface{}{
		"author":     reqDto.Author,
		"title":      reqDto.Title,
		"body":       reqDto.Body,
		"updated_at": time.Now(),
	}

	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	// Build Sql
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Update("articles").
		SetMap(values).
		Where(sq.Eq{"uuid": reqDto.Uuid}).
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sql, args...)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	article, err := c.findByUuid(ctx, reqDto.Uuid)
	if err != nil {
		return nil, err
	}

	// Push
	err = c.PushToQueue(ctx, articleUpdatedEvent, article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func (c *articleCommandPg) Delete(ctx context.Context, reqDto dto.RequestDeleteArticle) (*model.Article, error) {
	validate := validator.New()

	if err := validate.Struct(reqDto); err != nil {
		return nil, err
	}

	article, err := c.findByUuid(ctx, reqDto.Uuid)
	if err != nil {
		return nil, err
	}

	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	// Build Sql
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, err := psql.Delete("articles").
		Where(sq.Eq{"uuid": reqDto.Uuid}).
		ToSql()

	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, sql, args...)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	// Push
	err = c.PushToQueue(ctx, articleDeletedEvent, article)
	if err != nil {
		return nil, err
	}

	return article, nil
}
