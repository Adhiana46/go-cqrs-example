package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func (app *Config) GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("GET", os.Getenv("URL_QUERY_SVC")+"/articles", r.Body)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("Error calling GET /articles"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, response.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) GetSingleArticleHandler(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", os.Getenv("URL_QUERY_SVC"), "articles", uuid), r.Body)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("Error calling GET /articles/"+uuid))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, response.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) StoreArticleHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", os.Getenv("URL_COMMAND_SVC"), "articles"), r.Body)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("Error calling POST /articles"))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, response.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) UpdateArticleHandler(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s/%s", os.Getenv("URL_COMMAND_SVC"), "articles", uuid), r.Body)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("Error calling PUT /articles/"+uuid))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, response.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) DeleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s/%s", os.Getenv("URL_COMMAND_SVC"), "articles", uuid), r.Body)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code
	if response.StatusCode != http.StatusOK {
		app.errorJSON(w, errors.New("Error calling DELETE /articles/"+uuid))
		return
	}

	// create a variable we'll read response.Body into
	var jsonFromService jsonResponse

	// decode json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, response.StatusCode)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = jsonFromService.Message
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusOK, payload)
}
