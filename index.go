package quizzical

import (
	"auth"
	"controllers"
	asdatastore "datastore"
	"fmt"
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
	PARAM_CATEGORY       = "category"
	PARAM_QUESTION       = "question"
	PARAM_ANSWER         = "answer"
	PARAM_A              = "a"
	PARAM_B              = "b"
	PARAM_C              = "c"
	PARAM_D              = "d"
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

	/*
		http.Handle("/category", seshcookie.NewSessionHandler(&CategoryHandler{}, SESSION_KEY, nil))
		http.Handle("/question", seshcookie.NewSessionHandler(&QuestionHandler{}, SESSION_KEY, nil))
	*/

	m.Get("/", controllers.Index)
	m.Get("/logout", sessionauth.LoginRequired, controllers.Logout)
	m.Get("/admin", sessionauth.LoginRequired, func(w http.ResponseWriter, req *http.Request, r render.Render) {
		controllers.Admin(quizzicalService, w, req, r)
	})

	m.Post("/login", binding.Bind(models.User{}), controllers.Login)

	m.Get("/categories", quizzicalService.GetCategories)
	m.Get("/questions", quizzicalService.GetQuestions)

	http.Handle("/", m)
}

// GetAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() sessionauth.User {
	return &models.User{}
}

func (self *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue(PARAM_CATEGORY)

	if len(name) == 0 {
		http.Error(w, "Empty Category Name", http.StatusBadRequest)
		return
	}

	category := &models.Category{Name: name}
	err := categoryStore.Save(r, category)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Category '%s' added", category.Key)
}

func (self *QuestionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	question := r.FormValue(PARAM_QUESTION)
	answer := r.FormValue(PARAM_ANSWER)
	category := r.FormValue(PARAM_CATEGORY)
	a := r.FormValue(PARAM_A)
	b := r.FormValue(PARAM_B)
	c := r.FormValue(PARAM_C)
	d := r.FormValue(PARAM_D)

	if len(question) == 0 || len(answer) == 0 || len(a) == 0 || len(b) == 0 || len(c) == 0 || len(d) == 0 || len(category) == 0 {
		http.Error(w, "incomplete", http.StatusBadRequest)
	}

	q := &models.Question{Question: question, Answer: answer, Category: category, A: a, B: b, C: c, D: d}

	err := questionStore.Save(r, q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Question '%s' added", q.Key)
}
