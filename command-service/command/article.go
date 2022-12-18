package command

import (
	"context"
	"time"

	"github.com/Adhiana46/command-service/dto"
	"github.com/Adhiana46/command-service/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	articleCreatedEvent = ""
	articleUpdatedEvent = ""
	articleDeletedEvent = ""
)

type ArticleCommand interface {
	Store(ctx context.Context, reqDto dto.RequestStoreArticle) (*model.Article, error)
	Update(ctx context.Context, reqDto dto.RequestUpdateArticle) (*model.Article, error)
	Delete(ctx context.Context, reqDto dto.RequestDeleteArticle) (*model.Article, error)
}

type articleCommandPg struct {
	db         *sqlx.DB
	rabbitConn *amqp.Connection
}

func NewArticleCommandPg(db *sqlx.DB, rabbitConn *amqp.Connection) ArticleCommand {
	return &articleCommandPg{
		db:         db,
		rabbitConn: rabbitConn,
	}
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

	// TODO: publish event

	return c.findByUuid(ctx, values["uuid"].(string))
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
	// TODO: publish event

	return c.findByUuid(ctx, reqDto.Uuid)
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

	// TODO: publish event

	return article, nil
}
