package api

import (
	"datastore"
	"fmt"
	"auth"
	"github.com/martini-contrib/render"
	"github.com/dgrijalva/jwt-go"
	"models"
	"net/http"
	"utils"
	"bitbucket.org/waqqas-abdulkareem/asfour_toolkit/handlers"
)

const (
	ParameterLimit            = "limit"
	DefaultLimit              = 10
	DefaultLimitString        = "10"
	ParameterFormat           = "format"
	FormatXML                 = "xml"
	FormatJSON                = "json"
	ParameterCategory         = "category"
	HttpStatusValidationError = 422
)

func jwtHandlerFactory(handler func (http.ResponseWriter,*http.Request,*jwt.Token) error) *handlers.JWTRequestHandler{

	return handlers.QuickJWT(
		"token",
		[]byte(auth.JWTSecret),
		handler,
	)
}

func GetCategories(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render) {

	form := utils.FormHelper{Request:req}
	format := form.String(ParameterFormat,FormatXML)
	limit := form.Int(ParameterLimit,DefaultLimit)

	categories, err := dm.CategoryStore.GetAll(req, limit)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if format == FormatJSON {
		r.JSON(200, categories)
	} else {
		r.XML(200, models.Categories{Categories: categories})
	}
}

func GetQuestions(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render) {

	form := utils.FormHelper{Request: req}
	format := form.String(ParameterFormat,FormatXML)
	limit := form.Int(ParameterLimit,DefaultLimit)
	category := form.String(ParameterCategory,"")

	if len(category) == 0 {
		http.Error(w, fmt.Errorf("Parameter '%s' is missing.", ParameterCategory).Error(), HttpStatusValidationError)
		return
	}

	questions, err := dm.QuestionStore.Random(req, category, limit)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if format == FormatJSON {
		r.JSON(200, questions)
	} else {
		r.XML(200, models.Questions{Questions: questions, Category: category})
	}
}

func GetJWTCategories(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render){

	jwtHandlerFactory(func (w http.ResponseWriter,req *http.Request,token *jwt.Token) error{

		claims :=  utils.MapHelper{Map: token.Claims}
		format := claims.String(ParameterFormat,FormatXML)
		limit := claims.Int(ParameterLimit,DefaultLimit)

		categories, err := dm.CategoryStore.GetAll(req, limit)

		if err != nil { return err }

		if format == FormatJSON {
			r.JSON(200, categories)
		} else {
			r.XML(200, models.Categories{Categories: categories})
		}

		return nil
	
	
	}).ServeHTTP(w,req)
}

func GetJWTQuestions(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render) {

	jwtHandlerFactory(func (w http.ResponseWriter,req *http.Request,token *jwt.Token) error{

		claims := utils.MapHelper{Map: token.Claims}
		format := claims.String(ParameterFormat,FormatXML)
		limit := claims.Int(ParameterLimit,DefaultLimit)
		category := claims.String(ParameterCategory,"")

		if len(category) == 0 {
			http.Error(w, fmt.Errorf("Parameter '%s' is missing.", ParameterCategory).Error(), HttpStatusValidationError)
			return nil
		}

		questions, err := dm.QuestionStore.Random(req, category, limit)

		if err != nil { return err }

		if format == FormatJSON {
			r.JSON(200, questions)
		} else {
			r.XML(200, models.Questions{Questions: questions})
		}

		return nil	
	
	}).ServeHTTP(w,req)
}
