package api

import(
	"auth"
	"datastore"
	"net/http"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"appengine"
	"encoding/json"
)

type QuizzicalAPI struct{
	CategoryStore *datastore.CategoryStore
	QuestionStore *datastore.QuestionStore
	Consumer * jwt.Consumer
	Context appengine.Context
}

func (api * QuizzicalAPI) Authenticate(r * http.Request) (*jwt.Token,error){
	return api.Consumer.ValidateTokenFromRequestParameter(r,"token",[]byte(auth.Key))
}

func (api * QuizzicalAPI) Error(w http.ResponseWriter, err error,status int){

	
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)

	js, jsonerr := json.MarshalIndent(struct{Error string}{Error: err.Error()},"","  ")
	if jsonerr != nil{
		panic(jsonerr)
	}

	w.Write(js)
}

func (api * QuizzicalAPI) Success(w http.ResponseWriter, body interface{}){
 
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	js, err := json.MarshalIndent(body,"","  ")
	if err != nil {
		panic(err)
	}

	w.Write(js)
}

type APIHandler func(*http.Request, *jwt.Token, *QuizzicalAPI) (interface {}, error)

func (api *QuizzicalAPI) HandleWith(handler APIHandler) func(http.ResponseWriter, *http.Request){
	return func(w http.ResponseWriter, r * http.Request){

		token, err :=api.Authenticate(r)
		if err != nil {
			api.Error(w,err,http.StatusUnauthorized)
			return
		}

		api.Context = appengine.NewContext(r)

		result, err := handler(r,token,api)
		if err != nil {
			api.Error(w, err, http.StatusInternalServerError)
			return
		}

		api.Success(w, result)
	}
}