package main

import (
	"net/http"

	"github.com/Adhiana46/command-service/dto"
	"github.com/go-chi/chi/v5"
)

func (app *Config) StoreArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestDto dto.RequestStoreArticle
	_ = app.readJSON(w, r, &requestDto)

	article, err := app.cmdArticle.Store(ctx, requestDto)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Article Successfully Created",
		Data:    dto.ArticleToResponseDTO(article),
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Config) UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uuid := chi.URLParam(r, "uuid")

	var requestDto dto.RequestUpdateArticle
	_ = app.readJSON(w, r, &requestDto)
	requestDto.Uuid = uuid

	article, err := app.cmdArticle.Update(ctx, requestDto)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Article Successfully Updated",
		Data:    dto.ArticleToResponseDTO(article),
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *Config) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uuid := chi.URLParam(r, "uuid")

	requestDto := dto.RequestDeleteArticle{
		Uuid: uuid,
	}

	article, err := app.cmdArticle.Delete(ctx, requestDto)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Article Successfully Deleted",
		Data:    dto.ArticleToResponseDTO(article),
	}

	app.writeJSON(w, http.StatusOK, resp)
}
