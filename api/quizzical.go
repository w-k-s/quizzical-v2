package api

import(
	"auth"
	"datastore"
	"net/http"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"appengine"
)

type QuizzicalAPI struct{
	CategoryStore *datastore.CategoryStore
	QuestionStore *datastore.QuestionStore
	Consumer * jwt.Consumer
	Context * appengine.Context
}

func (api * QuizzicalAPI) Authenticate(r * http.Request) (jwt.Token,error){
	return consumer.ValidateTokenFromRequestParameter(r,"token",[]byte(auth.Key))
}

func (api * QuizzicalAPI) Error(w http.ResponseWriter, body interface{},status int){

}

func (api * QuizzicalAPI) Success(w http.ResponseWriter, body interface{}){
 
}

func (api *QuizzicalAPI) HandleWith(handler func(*http.Request, *QuizzicalAPI, *jwt.Token)) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, r * http.Request){

		token, err :=api.Authenticate(r)
		if err != nil {
			api.Error(w,err.Error(),http.StatusUnauthorized)
		}

		api.Context := appengine.NewContext(r)

		result, err := handler(r,api,token)
		if err != nil {
			api.Error(w, err.Error(), http.StatusInternalServerError)
		}

		api.Success(w, result)
	}
}