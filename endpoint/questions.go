package endpoint

import(
	"net/http"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"fmt"
	"api"
)


type QuestionEndpoint struct{}

var Question QuestionEndpoint;

func (endpoint *QuestionEndpoint) List(r * http.Request,token * jwt.Token, api * api.QuizzicalAPI) (interface{},error){

	category,valid := token.Claims[ParamNameCategory].(string);

	if !valid {
		return nil, fmt.Errorf("Required Parameter '%s' not present in token claims.",ParamNameCategory)
	}

	limit := int(token.Int32(ParamNameLimit,DefaultLimit))
	offset := int(token.Int32(ParamNameOffset,DefaultLimit))

	return api.QuestionStore.GetAll(api.Context, limit, offset, category)

}

func (endpoint *QuestionEndpoint) Post(r * http.Request,token * jwt.Token, api * api.QuizzicalAPI) (interface{},error){

	return nil,fmt.Errorf("No Implementation")
}