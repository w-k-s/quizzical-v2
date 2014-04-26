package quizzical

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"github.com/bpowers/seshcookie"
)

const (
	PARAM_USERNAME  = "username"
	PARAM_PASSWORD  = "password"
	STUPID_USERNAME = "waqqas"
	STUPID_PASSWORD = "CharlieAndTheChocolateFactory"
	PARAM_CATEGORY  = "category"
	PARAM_QUESTION  = "question"
	PARAM_ANSWER    = "answer"
	PARAM_A         = "a"
	PARAM_B         = "b"
	PARAM_C         = "c"
	PARAM_D         = "d"
	ENTITY_CATEGORY = "category"
	ENTITY_QUESTION = "question"
	KEY_AUTHENTICATED = "authenticated"
	//PARAM_STUPID_SECURITY = "stupid_security"
	SESSION_KEY       = "239ru238rhiou34hroi1uoi"
	//seperator          = string(os.PathSeparator)
	//QUESTIONS_FILE_EXT = ".xml"
	//QUESTIONS_DIR      = "questions"
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

func init() {
	http.HandleFunc("/", indexHandler)
	http.Handle("/auth", seshcookie.NewSessionHandler(&AuthHandler{},SESSION_KEY,nil))
	http.Handle("/admin", seshcookie.NewSessionHandler(&AdminHandler{},SESSION_KEY,nil))
	http.Handle("/category", seshcookie.NewSessionHandler(&CategoryHandler{},SESSION_KEY,nil))
	http.Handle("/question", seshcookie.NewSessionHandler(&QuestionHandler{},SESSION_KEY,nil))
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

func (self * AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
		//w.Header().Set("Status",string(http.StatusFound))
		//w.Header().Set("Location","/admin")
		authenticateAndRedirect(w,r,"/admin",http.StatusFound)
		return
	} else {
		http.Error(w, "Bad Credentials", http.StatusUnauthorized)
		return
	}
}

func (self * AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if redirectIfUnauthenticated(w,r,"/",http.StatusFound){
		return;
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

func (self * CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if redirectIfUnauthenticated(w,r,"/",http.StatusFound){
		return;
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
	//key := datastore.NewIncompleteKey(c, ENTITY_CATEGORY, nil)

	categories, err := listCategories(c)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responsdWithJSON(w, categories)
}

func (self * QuestionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	// path := QUESTIONS_DIR + seperator + category + QUESTIONS_FILE_EXT
	// bytes, err := ioutil.ReadFile(path)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "text/xml")
	// fmt.Fprint(w, string(bytes))
}

func listCategories(c appengine.Context) ([]Category, error) {

	query := datastore.NewQuery(ENTITY_CATEGORY)

	var categories []Category
	_, err := query.GetAll(c, &categories)

	if err != nil {
		return nil, err
	}

	// files, err := ioutil.ReadDir(QUESTIONS_DIR)
	// categories := make([]string, 0)
	// for _, file := range files {
	// 	fileName := file.Name()

	// 	indexExtension := strings.Index(fileName, QUESTIONS_FILE_EXT)

	// 	if indexExtension == -1 {
	// 		continue
	// 	}

	// 	categories = append(categories, file.Name()[:indexExtension])
	// }
	return categories, nil
}

func listQuestions(c appengine.Context, category string) ([]Question, error) {

	query := datastore.NewQuery(ENTITY_QUESTION).Filter(fmt.Sprintf("%s =", "Category"), category).Limit(20)

	var questions []Question
	_, err := query.GetAll(c, &questions)

	if err != nil {
		return nil, err
	}

	return questions, nil
}

func countQuestions(c appengine.Context, category string) (int,error){
	query := datastore.NewQuery(ENTITY_QUESTION).Filter(fmt.Sprintf("%s =","Category"), category)

	count,err := query.Count(c)

	if err != nil{
		return -1,err
	}

	return count,nil 
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


func authenticateAndRedirect(w http.ResponseWriter, r * http.Request, url string, status int) {
	session := seshcookie.Session.Get(r)
	session[KEY_AUTHENTICATED] = true
	http.Redirect(w, r, url, status)

}

func redirectIfUnauthenticated(w http.ResponseWriter,r * http.Request,url string, status int) bool{
	shouldRedirect := false
	session := seshcookie.Session.Get(r)

	if session == nil || session[KEY_AUTHENTICATED]==nil{
		shouldRedirect = true
	}else{
		authenticated := session[KEY_AUTHENTICATED].(bool)
		shouldRedirect = !authenticated
	}

	if shouldRedirect {
		http.Redirect(w,r,url,status)
	}

	return shouldRedirect
}

func responsdWithJSON(w http.ResponseWriter, v interface{}) {

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
