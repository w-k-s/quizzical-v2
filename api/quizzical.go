package api

import (
	"auth"
	"datastore"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"bitbucket.org/waqqas-abdulkareem/jwt-go"

	"appengine"
)

type QuizzicalAPI struct {
	CategoryStore *datastore.CategoryStore
	QuestionStore *datastore.QuestionStore
	Consumer      *jwt.Consumer
	Context       appengine.Context
}

type APIHandler func(*http.Request, *QuizzicalAPI) (interface{}, error)

func (api *QuizzicalAPI) HandleWith(handler APIHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		api.Context = appengine.NewContext(r)

		result, err := handler(r, api)
		if err != nil {
			api.Respond(w, r, NewError(err), http.StatusInternalServerError)
			return
		}

		api.Respond(w, r, result, http.StatusOK)

	}
}

func (api *QuizzicalAPI) AuthHandleWith(handler APIHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		_, err := api.Authenticate(r)
		if err != nil {
			api.Respond(w, r, NewError(err), http.StatusUnauthorized)
			return
		}

		api.HandleWith(handler)(w, r)
	}
}

func (api *QuizzicalAPI) Authenticate(r *http.Request) (*jwt.Token, error) {
	return api.Consumer.ValidateTokenFromRequestHeader(r, []byte(auth.Key))
}

func (api *QuizzicalAPI) Respond(w http.ResponseWriter, r *http.Request, body interface{}, status int) {

	var contentType string
	var content []byte
	var err error

	if r.Header.Get("Accept") == "application/xml" {

		contentType = "application/xml; charset=UTF-8"
		content, err = xml.Marshal(body)

	} else {

		contentType = "application/json; charset=UTF-8"
		content, err = json.MarshalIndent(body, "", "  ")
	}

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	w.Write(content)
}
