package quizzical

import (
	"api"
	"github.com/gorilla/mux"
	_"time"
	"net/http"
	"endpoint"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"datastore"
)


func init() {

	router := mux.NewRouter()
	api := setupQuizzicalAPI()

	apiSubrouter := router.PathPrefix("/api/v2/").Subrouter()

	apiSubrouter.HandleFunc("/categories",api.HandleWith(endpoint.Category.List)).
		Methods("GET")

	apiSubrouter.HandleFunc("/category",api.HandleWith(endpoint.Category.Post)).
		Methods("POST")

	apiSubrouter.HandleFunc("/questions",api.HandleWith(endpoint.Question.List)).
		Methods("GET")

	apiSubrouter.HandleFunc("/question",api.HandleWith(endpoint.Question.Post)).
		Methods("POST")


	http.Handle("/", router)
}

func setupQuizzicalAPI() *api.QuizzicalAPI{

	return &api.QuizzicalAPI{
		CategoryStore: &datastore.CategoryStore{},
		QuestionStore: &datastore.QuestionStore{},
		Consumer:	   setupConsumer(),
	}

}

func setupConsumer() *jwt.Consumer{

	consumer := jwt.NewConsumer("HS256")
	consumer.SetJTIRequired(true)
	consumer.SetExpirationTimeRequired(false)
	consumer.SetIssuedAtRequired(false)
	consumer.SetTokenLifespanInMinutesSinceIssue(2)

	return consumer
}