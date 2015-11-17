package quizzical

import (
	"api"
	"github.com/gorilla/mux"
	_"time"
	"net/http"
	"endpoint"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"datastore"
	"os"
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

	apiSubrouter.HandleFunc("/question",api.HandleWith(endpoint.Question.Delete)).
		Methods("DELETE")

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
	consumer.SetExpirationTimeRequired(!isDevelopmentMode())
	consumer.SetIssuedAtRequired(!isDevelopmentMode())
	consumer.SetTokenLifespanInMinutesSinceIssue(2)

	return consumer
}

func isDevelopmentMode() bool{

	var env []string = os.Environ()
	for i := 0; i < len(env); i++ {
		if env[i] == "RUN_WITH_DEVAPPSERVER=1" {
			return true
		}
	}

	return false
}