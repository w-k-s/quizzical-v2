package quizzical

import (
	"api"
	"auth"
	"controllers"
	"datastore"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"models"
	"net/http"
)

const (
	SessionName          = "com.appspot.asfour-quizzical.session"
	SessionRedirectUrl   = "/"
	SessionRedirectParam = "forward"
)

type Encoder interface {
	Encode(v interface{}) error
}

var categoryStore *datastore.CategoryStore
var questionStore *datastore.QuestionStore

func init() {

	registerModelsForTemplating()

	categoryStore = &datastore.CategoryStore{}
	questionStore = &datastore.QuestionStore{}
	dataManager := &datastore.Manager{CategoryStore: categoryStore, QuestionStore: questionStore}

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		IndentJSON: true,
		IndentXML:  true,
		Charset:    "UTF-8",
	}))

	store := sessions.NewCookieStore([]byte(auth.SessionAuthenticationKey))
	m.Use(sessions.Sessions(SessionName, store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
	sessionauth.RedirectUrl = SessionRedirectUrl
	sessionauth.RedirectParam = SessionRedirectParam

	//Allow Martini to inject the datastore manager as a service.
	m.Map(dataManager)

	m.Get("/", controllers.GetIndex)
	m.Get("/logout", sessionauth.LoginRequired, controllers.GetLogout)
	m.Get("/admin", sessionauth.LoginRequired, controllers.GetAdmin)

	m.Post("/login", binding.Bind(models.User{}), controllers.PostLogin)
	m.Post("/category", sessionauth.LoginRequired, binding.Bind(models.Category{}), controllers.PostCategory)
	m.Post("/question", sessionauth.LoginRequired, binding.Bind(models.Question{}), controllers.PostQuestion)

	m.Get("/api/categories", api.GetJWTCategories)
	m.Get("/api/questions", api.GetJWTQuestions)

	//-------------------------------------------------------------------------//

	consumer := jwt.NewConsumer("HS256")
	consumer.SetJTIRequired(true)
	consumer.SetExpirationTimeRequired(true)
	consumer.SetIssuedAtRequired(true)
	consumer.SetTokenLifespanInMinutesSinceIssue(2 * 60 * 1000)

	api2 := api.QuizzicalApi{
		Consumer: consumer,
		DB:       dataManager,
		ResponseFormatter: func(r *http.Request, w http.ResponseWriter, response interface{}, err error) {

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var contentType string
			var encoder Encoder

			if r.FormValue("format") == "json" {
				contentType = "application/json; charset=UTF-8"
				encoder = json.NewEncoder(w)
			} else {
				contentType = "text/xml"
				encoder = xml.NewEncoder(w)
			}

			w.Header().Set("Content-Type", contentType)
			w.WriteHeader(http.StatusOK)
			if err = encoder.Encode(response); err != nil {
				panic(err)
			}
		},
	}

	m.Get("/api/v2/categories", api2.Categories)
	m.Get("/api/v2/questions", api2.Questions)
	m.Post("/api/v2/category", api2.PostCategory)

	//-------------------------------------------------------------------------//

	http.Handle("/", m)
}

// GetAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() sessionauth.User {
	return &models.User{}
}

func registerModelsForTemplating() {

	gob.Register(models.Category{})
	gob.Register(models.Question{})

}
