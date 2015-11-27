package endpoint

import (
	"api"
	"encoding/json"
	_ "fmt"
	"models"
	"net/http"
	"utils"
	"io/ioutil"
)

type CategoryEndpoint struct{}

var Category CategoryEndpoint

func (endpoint *CategoryEndpoint) List(r *http.Request, quizApi *api.QuizzicalAPI) (interface{}, error) {
	
	categories, err := quizApi.CategoryStore.GetAll(quizApi.Context, utils.FormUInt(r,ParamNamePageSize,DefaultPageSize))

	return api.Response{Data: categories}, err
}

func (endpoint *CategoryEndpoint) Post(r *http.Request, quizApi *api.QuizzicalAPI) (interface{}, error) {


	body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return nil, err;
    }

	var category models.Category
	err = json.Unmarshal(body, &category)
    if err != nil {
        return nil, err;
    }

	err = category.Validate()

	if err != nil {
		return nil, err
	}

	err = quizApi.CategoryStore.Save(quizApi.Context, &category)

	return api.Response{Data: category}, err
}
