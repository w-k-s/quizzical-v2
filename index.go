package quizzical

import (
	"api"
	"auth"
	"github.com/gorilla/mux"
	"time"
	"net/http"
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

func setupQuizzicalAPI() *QuizzicalAPI{

	return &QuizzicalApi{
		CategoryStore: &datastore.CategoryStore{},
		QuestionStore: &datastore.QuestionStore{},
		Consumer:	   setupConsumer(),
	}

}

func setupConsumer() *jwt.Consumer{

	consumer := jwt.NewConsumer("HS256")
	consumer.SetJTIRequired(true)
	consumer.SetExpirationTimeRequired(true)
	consumer.SetIssuedAtRequired(true)
	consumer.SetTokenLifespanInMinutesSinceIssue(2 * time.Minute)

	return consumer
}