package endpoint

import (
	"api"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"encoding/json"
	_ "fmt"
	"models"
	"net/http"
)

type CategoryEndpoint struct{}

var Category CategoryEndpoint

func (endpoint *CategoryEndpoint) List(r *http.Request, token *jwt.Token, quizApi *api.QuizzicalAPI) (interface{}, error) {

	categories, err := quizApi.CategoryStore.GetAll(quizApi.Context, int(token.Int32(ParamNameLimit, DefaultLimit)))

	return api.Response{Data: categories}, err
}

func (endpoint *CategoryEndpoint) Post(r *http.Request, token *jwt.Token, quizApi *api.QuizzicalAPI) (interface{}, error) {

	var category models.Category

	jsonCategory, err := json.Marshal(token.Claims[ParamNameCategory])

	if jsonCategory != nil {

		_ = json.Unmarshal(jsonCategory, &category)

		err = category.Validate()

		if err == nil {

			err = quizApi.CategoryStore.Save(quizApi.Context, &category)

		}

	}

	return api.Response{Data: category}, err
}
