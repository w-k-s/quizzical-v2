package endpoint

import (
	"api"
	"encoding/json"
	"fmt"
	"models"
	"utils"
	"net/http"
	"io/ioutil"
)

type QuestionEndpoint struct{}

var Question QuestionEndpoint

func (endpoint *QuestionEndpoint) List(r *http.Request, quizApi *api.QuizzicalAPI) (interface{}, error) {

	category := r.FormValue(ParamNameCategory)

	if len(category) == 0 {
		return nil, fmt.Errorf("Required Parameter '%s' not present in query.", ParamNameCategory)
	}

	pageSize := utils.FormUInt(r,ParamNamePageSize,DefaultPageSize)
	pageNumber := utils.FormUInt(r,ParamNamePageNumber,DefaultPageNumber)

	result, err := quizApi.QuestionStore.GetAll(quizApi.Context, pageSize, pageNumber, category)

	return api.NewPaginatedResponse(result.Data,pageSize,pageNumber,result.TotalCount), err
}

func (endpoint *QuestionEndpoint) Post(r *http.Request, quizApi *api.QuizzicalAPI) (interface{}, error) {

	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return nil, err;
    }

    var question models.Question
    
	err = json.Unmarshal(body, &question)
    if err != nil {
        return nil, err;
    }
	err = question.Validate()

	if err != nil {
		return nil, err
	}

	err = quizApi.QuestionStore.Save(quizApi.Context, &question)

	return api.Response{Data: question}, err
}

func (endpoint *QuestionEndpoint) PostMulti(r *http.Request, quizApi *api.QuizzicalAPI) (interface{}, error) {

	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return nil, err;
    }

    var questions []models.Question
    
	err = json.Unmarshal(body, &questions)
    if err != nil {
        return nil, err;
    }

    for _,question := range questions{
    	err  = question.Validate()
    	if  err != nil{
    		return nil,err
    	}
    }

	err = quizApi.QuestionStore.SaveMulti(quizApi.Context, questions)

	return api.Response{Data: questions}, err
}

func (endpoint *QuestionEndpoint) Delete(r *http.Request, quizApi *api.QuizzicalAPI) (interface{}, error) {

	key := r.FormValue(ParamNameKey)

	if len(key) == 0 {
		return nil, fmt.Errorf("Required Parameter '%s' not present in query.", ParamNameKey)
	}

	err := quizApi.QuestionStore.Delete(quizApi.Context, key)

	return api.Response{}, err
}
