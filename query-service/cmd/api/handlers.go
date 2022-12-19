package main

import (
	"net/http"
	"strconv"

	"github.com/Adhiana46/query-service/dto"
	"github.com/go-chi/chi/v5"
)

func (app *Config) GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page == 0 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit == 0 {
		limit = 25
	}

	q := r.URL.Query().Get("q")
	author := r.URL.Query().Get("author")

	requestDto := dto.RequestListArticle{
		Page:   page,
		Limit:  limit,
		Query:  q,
		Author: author,
	}

	articles, err := app.queryArticle.GetList(ctx, requestDto)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Succesfully Get List of Articles",
		Data:    dto.ArticlesToResponseDtos(articles),
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Config) GetSingleArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uuid := chi.URLParam(r, "uuid")

	requestDto := dto.RequestSingleArticle{
		Uuid: uuid,
	}

	article, err := app.queryArticle.GetSingle(ctx, requestDto)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Sucessfully Get Article",
		Data:    dto.ArticleToResponseDTO(article),
	}

	app.writeJSON(w, http.StatusOK, resp)
}
