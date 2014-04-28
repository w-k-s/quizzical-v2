package quizzical

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/bpowers/seshcookie"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"text/template"
)

const (
	PARAM_USERNAME    = "username"
	PARAM_PASSWORD    = "password"
	STUPID_USERNAME   = "waqqas"
	STUPID_PASSWORD   = "CharlieAndTheChocolateFactory"
	PARAM_CATEGORY    = "category"
	PARAM_QUESTION    = "question"
	PARAM_ANSWER      = "answer"
	PARAM_A           = "a"
	PARAM_B           = "b"
	PARAM_C           = "c"
	PARAM_D           = "d"
	PARAM_OLD_CATEGORY="old_category"
	PARAM_OLD_QUESTIOM= "old_question"
	ENTITY_CATEGORY   = "category"
	ENTITY_QUESTION   = "question"
	KEY_AUTHENTICATED = "authenticated"
	SESSION_KEY       = "239ru238rhiou34hroi1uoi"
	NUM_QUESTIONS     = 20
)

type Category struct {
	Name string
}

type Question struct {
	XMLName  xml.Name `xml:"question"`
	Question string   `xml:"ask"`
	Answer   string   `xml:"correct,attr"`
	Category string   `xml:"-"`
	A        string   `xml:"A"`
	B        string   `xml:"B"`
	C        string   `xml:"C"`
	D        string   `xml:"D"`
}

type Questions struct {
	XMLName   xml.Name `xml:"questions"`
	Category  string   `xml:"title,attr"`
	Questions []Question
}

type AuthHandler struct{}
type AdminHandler struct{}
type CategoryHandler struct{}
type QuestionHandler struct{}
type EditQuestionHandler struct{}

func init() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/auth", seshcookie.NewSessionHandler(&AuthHandler{}, SESSION_KEY, nil))
	http.Handle("/admin", seshcookie.NewSessionHandler(&AdminHandler{}, SESSION_KEY, nil))
	http.Handle("/category", seshcookie.NewSessionHandler(&CategoryHandler{}, SESSION_KEY, nil))
	http.Handle("/question", seshcookie.NewSessionHandler(&QuestionHandler{}, SESSION_KEY, nil))
	http.Handle("/question/edit", seshcookie.NewSessionHandler(&EditQuestionHandler{},SESSION_KEY,nil))
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/questions", questionsHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	html, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(html))
}

func (self *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue(PARAM_USERNAME)
	password := r.FormValue(PARAM_PASSWORD)

	if len(username) == 0 {
		http.Error(w, "Username not provided", http.StatusBadRequest)
		return
	} else if len(password) == 0 {
		http.Error(w, "Password not provided", http.StatusBadRequest)
		return
	}

	if strings.EqualFold(username, STUPID_USERNAME) && strings.EqualFold(password, STUPID_PASSWORD) {
		authenticateAndRedirect(w, r, "/admin", http.StatusFound)
		return
	} else {
		http.Error(w, "Bad Credentials", http.StatusUnauthorized)
		return
	}
}

func (self *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if redirectIfUnauthenticated(w, r, "/", http.StatusFound) {
		return
	}

	html, err := ioutil.ReadFile("admin.html")
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

	category := &Category{Name: name}
	c := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(c, ENTITY_CATEGORY, nil)
	_, err := datastore.Put(c, key, category)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Category '%s' added", name)
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	categories, err := listCategories(c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, categories)
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

	q := &Question{Question: question, Answer: answer, Category: category, A: a, B: b, C: c, D: d}
	ctx := appengine.NewContext(r)

	key := datastore.NewIncompleteKey(ctx, ENTITY_QUESTION, nil)
	_, err := datastore.Put(ctx, key, q)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Question '%s' added", q.Question)
}

func (self * EditQuestionHandler) ServeHTTP(w http.ResponseWriter, r * http.Request){

	if redirectIfUnauthenticated(w,r,"/",http.StatusFound) {
		return;
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
	
	if len(oldQuestion) == 0 || len(oldCategory) == 0  {
		http.Error(w, "incomplete", http.StatusBadRequest)
		return;
	}
	
	context := appengine.NewContext(r)
	query :=datastore.NewQuery(ENTITY_QUESTION).Filter("Question = ",oldQuestion).Filter("Category = ",oldCategory)
	numResults,err := query.Count(context)

	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	if numResults > 1{
		respondWithText(w,"Too many results")
		return;
	}else if numResults < 1{
		respondWithText(w, "No matches")
		return;
	}

	var questions []*Question
	keys,err := query.GetAll(context,&questions)

	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}


	if len(question) == 0{
		question = questions[0].Question
	}else if len(answer) == 0{
		answer = questions[0].Answer
	}else if len(category) == 0{
		category = questions[0].Category
	}else if len(a) == 0{
		a = questions[0].A
	}else if len(b) == 0{
		b = questions[0].B
	}else if len(c) == 0{
		c = questions[0].C
	}else if len(c) == 0{
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

	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}else{
		respondWithText(w,fmt.Sprintf("Question '%v' updated to '%v'.",oldQuestion,question))
	}
}

func questionsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue(PARAM_CATEGORY)

	if len(category) == 0 {
		http.Error(w, "empty category", http.StatusBadRequest)
		return
	}

	context := appengine.NewContext(r)
	exists, err := categoryExists(context, category)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, "Category Not Found", http.StatusNotFound)
		return
	}

	questions, err := listQuestions(context, category)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	questionsList := &Questions{Questions: questions, Category: category}

	respondWithXML(w, questionsList)
}

func listCategories(c appengine.Context) ([]Category, error) {

	query := datastore.NewQuery(ENTITY_CATEGORY)

	var categories []Category
	_, err := query.GetAll(c, &categories)

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func listQuestions(c appengine.Context, category string) ([]Question, error) {

	count, err := countQuestions(c, category)

	var questions []Question
	query := datastore.NewQuery(ENTITY_QUESTION).Filter(fmt.Sprintf("%s =", "Category"), category)

	if err != nil && count > NUM_QUESTIONS {
		randomLimit := count - NUM_QUESTIONS
		start := rand.Int63n(int64(randomLimit))

		query.Offset(int(start))
	}

	_, err = query.Limit(NUM_QUESTIONS).GetAll(c, &questions)

	if err != nil {
		return nil, err
	}

	return questions, nil
}

func countQuestions(c appengine.Context, category string) (int, error) {
	query := datastore.NewQuery(ENTITY_QUESTION).Filter(fmt.Sprintf("%s =", "Category"), category)

	count, err := query.Count(c)

	if err != nil {
		return -1, err
	}

	return count, nil
}

func categoryExists(c appengine.Context, category string) (bool, error) {
	categories, err := listCategories(c)

	if err != nil {
		return false, err
	}

	for _, aCategory := range categories {
		if strings.EqualFold(aCategory.Name, category) {
			return true, nil
		}
	}

	return false, nil
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

func respondWithText(w http.ResponseWriter,text string){
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, text)
}

func respondWithJSON(w http.ResponseWriter, v interface{}) {

	json, err := json.MarshalIndent(v, "", "    ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(json))

}

func respondWithXML(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "text/xml")
	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
