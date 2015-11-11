package endpoint

import(
	"net/http"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"fmt"
	"api"
	"models"
	"encoding/json"
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

	questions, err :=  api.QuestionStore.GetAll(api.Context, limit, offset, category)

	return Response{Data: questions},err
}

func (endpoint *QuestionEndpoint) Post(r * http.Request,token * jwt.Token, api * api.QuizzicalAPI) (interface{},error){

	var question models.Question

	jsonQuestion,err := json.Marshal(token.Claims[ParamNameQuestion])
	
	if jsonQuestion != nil {
		
		_ = json.Unmarshal(jsonQuestion,&question)
		
		err = question.Validate()

		if err == nil {
			
			err = api.QuestionStore.Save(api.Context,&question)

		}else{
			err = ErrorWrapper{Message: err}
		}
	}

	return Response{Data: question},err
}