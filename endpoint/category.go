package endpoint

import(
	"net/http"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"api"
	_"fmt"
	"models"
	"encoding/json"

)

type CategoryEndpoint struct{}

var Category CategoryEndpoint;

func (endpoint *CategoryEndpoint) List(r * http.Request,token * jwt.Token, api * api.QuizzicalAPI) (interface{},error){
	
	categories,err := api.CategoryStore.GetAll(api.Context, int(token.Int32(ParamNameLimit,DefaultLimit)))

	return Response{Data: categories},err
}

func (endpoint *CategoryEndpoint) Post(r * http.Request,token * jwt.Token, api * api.QuizzicalAPI) (interface{},error){

	var category models.Category

	jsonCategory,err := json.Marshal(token.Claims[ParamNameCategory])
	
	if jsonCategory != nil {
		
		_ = json.Unmarshal(jsonCategory,&category)
		
		err = category.Validate()

		if err == nil {
			
			err = api.CategoryStore.Save(api.Context,&category)	
		
		}else{
			err = ErrorWrapper{OriginalError: err}
		}
		
	}

	return Response{Data: category},err
}