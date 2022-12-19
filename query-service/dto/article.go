package dto

import (
	"time"

	"github.com/Adhiana46/query-service/model"
)

type ResponseArticle struct {
	Uuid      string    `json:"uuid"`
	Author    string    `json:"author"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RequestSingleArticle struct {
	Uuid string `validate:"required"`
}

type RequestListArticle struct {
	Page   int    `json:"page" validate:""`
	Limit  int    `json:"limit" validate:""`
	Query  string `json:"query" validate:""`
	Author string `json:"author" validate:""`
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

func ArticlesToResponseDtos(articles []*model.Article) []*ResponseArticle {
	result := []*ResponseArticle{}
	for _, article := range articles {
		result = append(result, ArticleToResponseDTO(article))
	}
	return result
}
