package endpoint

import(
	"net/http"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"models"
)

const DefaultLimit = 10

type CategoryEndpoint struct{}

var Category CategoryEndpoint;

func (endpoint *CategoryEndpoint) List(r * http.Request,token * jwt.Token, api * QuizzicalAPI) (interface{},error){

	return api.CategoryStore.GetAll(api.Context, token.Int32("limit",DefaultLimit))

}

func (endpoint *CategoryEndpoint) Post(r * http.Request,token * jwt.Token, api * QuizzicalAPI) (interface{},error){


}