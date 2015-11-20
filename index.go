package quizzical

import (
	"api"
	"bitbucket.org/waqqas-abdulkareem/jwt-go"
	"datastore"
	"endpoint"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	_ "time"
)

func init() {

	router := mux.NewRouter()
	api := setupQuizzicalAPI()

	apiSubrouter := router.PathPrefix("/api/v2/").Subrouter()

	apiSubrouter.HandleFunc("/auth", api.HandleWith(endpoint.Auth.IssueToken)).
		Methods("POST").
		Schemes("https")

	apiSubrouter.HandleFunc("/categories", api.AuthHandleWith(endpoint.Category.List)).
		Methods("GET")

	apiSubrouter.HandleFunc("/category", api.AuthHandleWith(endpoint.Category.Post)).
		Methods("POST")

	apiSubrouter.HandleFunc("/questions", api.AuthHandleWith(endpoint.Question.List)).
		Methods("GET")

	apiSubrouter.HandleFunc("/question", api.AuthHandleWith(endpoint.Question.Post)).
		Methods("POST")

	apiSubrouter.HandleFunc("/question", api.AuthHandleWith(endpoint.Question.Delete)).
		Methods("DELETE")


	http.Handle("/", router)
}

func setupQuizzicalAPI() *api.QuizzicalAPI {

	return &api.QuizzicalAPI{
		CategoryStore: &datastore.CategoryStore{},
		QuestionStore: &datastore.QuestionStore{},
		Consumer:      setupConsumer(),
	}

}

func setupConsumer() *jwt.Consumer {

	consumer := jwt.NewConsumer("HS256")
	consumer.SetJTIRequired(true)
	consumer.SetExpirationTimeRequired(!isDevelopmentMode())
	consumer.SetIssuedAtRequired(!isDevelopmentMode())
	consumer.SetTokenLifespanInMinutesSinceIssue(60)

	return consumer
}

func isDevelopmentMode() bool {

	var env []string = os.Environ()
	for i := 0; i < len(env); i++ {
		if env[i] == "RUN_WITH_DEVAPPSERVER=1" {
			return true
		}
	}

	return false
}
