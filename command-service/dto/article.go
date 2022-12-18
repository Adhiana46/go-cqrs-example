package dto

import (
	"time"

	"github.com/Adhiana46/command-service/model"
)

type ResponseArticle struct {
	Uuid      string    `json:"uuid"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequestStoreArticle struct {
	Author string `json:"author" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

type RequestUpdateArticle struct {
	Uuid   string `validate:"required"`
	Author string `json:"author" validate:"required"`
	Title  string `json:"title" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

type RequestDeleteArticle struct {
	Uuid string `validate:"required"`
}

func ArticleToResponseDTO(article *model.Article) *ResponseArticle {
	return &ResponseArticle{
		Uuid:      article.Uuid,
		Author:    article.Author,
		Title:     article.Title,
		Body:      article.Body,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
}
