package endpoint

import (
	"api"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
)

type QuestionEndpoint struct{}

var Question QuestionEndpoint

func (endpoint *QuestionEndpoint) List(r *http.Request, token *jwt.Token, quizApi *api.QuizzicalAPI) (interface{}, error) {

	category, valid := token.Claims[ParamNameCategory].(string)

	if !valid {
		return nil, fmt.Errorf("Required Parameter '%s' not present in token claims.", ParamNameCategory)
	}

	limit := int(token.GetInt32(ParamNameLimit, DefaultLimit))
	offset := int(token.GetInt32(ParamNameOffset, DefaultOffset))

	questions, err := quizApi.QuestionStore.GetAll(quizApi.Context, limit, offset, category)

	return api.Response{Data: questions}, err
}

func (endpoint *QuestionEndpoint) Post(r *http.Request, token *jwt.Token, quizApi *api.QuizzicalAPI) (interface{}, error) {

	var question models.Question

	jsonQuestion, err := json.Marshal(token.Claims[ParamNameQuestion])

	if jsonQuestion != nil {

		_ = json.Unmarshal(jsonQuestion, &question)

		err = question.Validate()

		if err == nil {

			err = quizApi.QuestionStore.Save(quizApi.Context, &question)

		}
	}

	return api.Response{Data: question}, err
}

func (endpoint *QuestionEndpoint) Delete(r *http.Request, token *jwt.Token, quizApi *api.QuizzicalAPI) (interface{}, error) {

	key, present := token.Claims[ParamNameKey].(string)

	if !present {
		return nil, fmt.Errorf("Required Parameter '%s' not present in token claims.", ParamNameKey)
	}

	err := quizApi.QuestionStore.Delete(quizApi.Context, key)

	return api.Response{}, err
}
