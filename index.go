package quizzical

import (
	"auth"
	"controllers"
	asdatastore "datastore"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"models"
	"net/http"
	"services"
)

const (
	SessionName          = "com.appspot.asfour-quizzical.session"
	SessionRedirectUrl   = "/"
	SessionRedirectParam = "forward"
)

var categoryStore *asdatastore.CategoryStore
var questionStore *asdatastore.QuestionStore

type AuthHandler struct{}
type AdminHandler struct{}
type CategoryHandler struct{}
type QuestionHandler struct{}
type EditQuestionHandler struct{}

func init() {

	categoryStore = &asdatastore.CategoryStore{}
	questionStore = &asdatastore.QuestionStore{}
	quizzicalService := &services.QuizzicalService{CategoryStore: categoryStore, QuestionStore: questionStore}

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

	m.Get("/", controllers.GetIndex)
	m.Get("/logout", sessionauth.LoginRequired, controllers.GetLogout)
	m.Get("/admin", sessionauth.LoginRequired, func(w http.ResponseWriter, req *http.Request, r render.Render) {
		controllers.GetAdmin(quizzicalService, w, req, r)
	})

	m.Post("/login", binding.Bind(models.User{}), controllers.PostLogin)
	m.Post("/category", binding.Bind(models.Category{}), func(w http.ResponseWriter, req *http.Request, postedCategory models.Category, r render.Render) {
		controllers.PostCategory(quizzicalService, postedCategory, w, req, r)
	})
	m.Post("/question", binding.Bind(models.Question{}), func(w http.ResponseWriter, req *http.Request, postedQuestion models.Question, r render.Render) {
		controllers.PostQuestion(quizzicalService, postedQuestion, w, req, r)
	})

	m.Get("/categories", quizzicalService.GetCategories)
	m.Get("/questions", quizzicalService.GetQuestions)

	http.Handle("/", m)
}

// GetAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() sessionauth.User {
	return &models.User{}
}
