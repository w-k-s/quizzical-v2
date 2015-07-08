package quizzical

import (
	"appengine"
	"appengine/datastore"
	"auth"
	"bytes"
	"controllers"
	asdatastore "datastore"
	"fmt"
	"github.com/bpowers/seshcookie"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"io/ioutil"
	"models"
	"net/http"
	"services"
	"text/template"
	"utils"
)

const (
	PARAM_USERNAME       = "username"
	PARAM_PASSWORD       = "password"
	STUPID_USERNAME      = "waqqas"
	STUPID_PASSWORD      = "CharlieAndTheChocolateFactory"
	PARAM_CATEGORY       = "category"
	PARAM_QUESTION       = "question"
	PARAM_ANSWER         = "answer"
	PARAM_A              = "a"
	PARAM_B              = "b"
	PARAM_C              = "c"
	PARAM_D              = "d"
	PARAM_OLD_CATEGORY   = "old_category"
	PARAM_OLD_QUESTIOM   = "old_question"
	ENTITY_CATEGORY      = "category"
	ENTITY_QUESTION      = "question"
	KEY_AUTHENTICATED    = "authenticated"
	SESSION_KEY          = "239ru238rhiou34hroi1uoi"
	NUM_QUESTIONS        = 20
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

	store := sessions.NewCookieStore([]byte(auth.SessionAuthenticationKey), []byte(auth.SessionEncryptionKey))
	m.Use(sessions.Sessions(SessionName, store))
	m.Use(sessionauth.SessionUser(GenerateAnonymousUser))
	sessionauth.RedirectUrl = SessionRedirectUrl
	sessionauth.RedirectParam = SessionRedirectParam

	/*http.HandleFunc("/", indexHandler)
	http.Handle("/auth", seshcookie.NewSessionHandler(&AuthHandler{}, SESSION_KEY, nil))
	http.Handle("/admin", seshcookie.NewSessionHandler(&AdminHandler{}, SESSION_KEY, nil))
	http.Handle("/category", seshcookie.NewSessionHandler(&CategoryHandler{}, SESSION_KEY, nil))
	http.Handle("/question", seshcookie.NewSessionHandler(&QuestionHandler{}, SESSION_KEY, nil))
	http.Handle("/question/edit", seshcookie.NewSessionHandler(&EditQuestionHandler{},SESSION_KEY,nil))
	*/

	m.Get("/", controllers.Index)
	m.Get("/categories", quizzicalService.GetCategories)
	m.Get("/questions", quizzicalService.GetQuestions)

	m.Post("/login", binding.Bind(models.User{}), controllers.Login)

	http.Handle("/", m)
	http.Handle("/categories", m)
	http.Handle("/questions", m)
}

// GetAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() sessionauth.User {
	return &models.User{}
}

func (self *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if redirectIfUnauthenticated(w, r, "/", http.StatusFound) {
		return
	}

	html, err := ioutil.ReadFile("templates/admin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.New("Admin")
	t, err = t.Parse(string(html))

	c := appengine.NewContext(r)
	categories, err := listCategories(c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var categoriesOptions bytes.Buffer
	for _, category := range categories {
		categoriesOptions.WriteString(fmt.Sprintf("<option value='%s'>%s</option>", category.Name, category.Name))
	}

	t.Execute(w, categoriesOptions.String())
}

func (self *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if redirectIfUnauthenticated(w, r, "/", http.StatusFound) {
		return
	}

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

func (self *EditQuestionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if redirectIfUnauthenticated(w, r, "/", http.StatusFound) {
		return
	}

	oldQuestion := r.FormValue(PARAM_OLD_QUESTIOM)
	oldCategory := r.FormValue(PARAM_OLD_CATEGORY)
	question := r.FormValue(PARAM_QUESTION)
	answer := r.FormValue(PARAM_ANSWER)
	category := r.FormValue(PARAM_CATEGORY)
	a := r.FormValue(PARAM_A)
	b := r.FormValue(PARAM_B)
	c := r.FormValue(PARAM_C)
	d := r.FormValue(PARAM_D)

	if len(oldQuestion) == 0 || len(oldCategory) == 0 {
		http.Error(w, "incomplete", http.StatusBadRequest)
		return
	}

	context := appengine.NewContext(r)
	query := datastore.NewQuery(ENTITY_QUESTION).Filter("Question = ", oldQuestion).Filter("Category = ", oldCategory)
	numResults, err := query.Count(context)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if numResults > 1 {
		utils.RespondWithText(w, "Too many results")
		return
	} else if numResults < 1 {
		utils.RespondWithText(w, "No matches")
		return
	}

	var questions []*models.Question
	keys, err := query.GetAll(context, &questions)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(question) == 0 {
		question = questions[0].Question
	} else if len(answer) == 0 {
		answer = questions[0].Answer
	} else if len(category) == 0 {
		category = questions[0].Category
	} else if len(a) == 0 {
		a = questions[0].A
	} else if len(b) == 0 {
		b = questions[0].B
	} else if len(c) == 0 {
		c = questions[0].C
	} else if len(c) == 0 {
		d = questions[0].D
	}

	questions[0].Question = question
	questions[0].Answer = answer
	questions[0].Category = category
	questions[0].A = a
	questions[0].B = b
	questions[0].C = c
	questions[0].D = d

	_, err = datastore.Put(context, keys[0], questions[0])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		utils.RespondWithText(w, fmt.Sprintf("Question '%v' updated to '%v'.", oldQuestion, question))
	}
}

func listCategories(c appengine.Context) ([]models.Category, error) {

	query := datastore.NewQuery(ENTITY_CATEGORY)

	var categories []models.Category
	_, err := query.GetAll(c, &categories)

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func authenticateAndRedirect(w http.ResponseWriter, r *http.Request, url string, status int) {
	session := seshcookie.Session.Get(r)
	session[KEY_AUTHENTICATED] = true
	http.Redirect(w, r, url, status)

}

func redirectIfUnauthenticated(w http.ResponseWriter, r *http.Request, url string, status int) bool {
	shouldRedirect := false
	session := seshcookie.Session.Get(r)

	if session == nil || session[KEY_AUTHENTICATED] == nil {
		shouldRedirect = true
	} else {
		authenticated := session[KEY_AUTHENTICATED].(bool)
		shouldRedirect = !authenticated
	}

	if shouldRedirect {
		http.Redirect(w, r, url, status)
	}

	return shouldRedirect
}
