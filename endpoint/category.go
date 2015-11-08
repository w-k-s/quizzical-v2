package endpoint

import(
	"net/http"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"api"
	"fmt"
)

type CategoryEndpoint struct{}

var Category CategoryEndpoint;

func (endpoint *CategoryEndpoint) List(r * http.Request,token * jwt.Token, api * api.QuizzicalAPI) (interface{},error){

	return api.CategoryStore.GetAll(api.Context, int(token.Int32("limit",DefaultLimit)))

}

func (endpoint *CategoryEndpoint) Post(r * http.Request,token * jwt.Token, api * api.QuizzicalAPI) (interface{},error){

	return nil,fmt.Errorf("No Implementation")
}