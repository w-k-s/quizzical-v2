package api

import (
	"auth"
	JWT "bitbucket.org/waqqas-abdulkareem/asfour_toolkit/handlers/jwt"
	"datastore"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/martini-contrib/render"
	"models"
	"net/http"
	"utils"
)

const (
	ParameterLimit            = "limit"
	ParameterOffset						= "offset"
	DefaultLimit              = 10
	DefaultOffset							= 0
	DefaultLimitString        = "10"
	ParameterFormat           = "format"
	FormatXML                 = "xml"
	FormatJSON                = "json"
	ParameterCategory         = "category"
	HttpStatusValidationError = 422
)

func jwtHandlerFactory(handler func(http.ResponseWriter, *http.Request, *jwt.Token) error) *JWT.Handler {

	return JWT.NewHandler(
		"token",
		[]byte(auth.JWTSecret),
		handler,
	)
}

func GetJWTCategories(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render) {

	jwtHandlerFactory(func(w http.ResponseWriter, req *http.Request, token *jwt.Token) error {

		claims := utils.MapHelper{Map: token.Claims}
		format := claims.String(ParameterFormat, FormatXML)
		limit := claims.Int(ParameterLimit, DefaultLimit)

		categories, err := dm.CategoryStore.GetAll(req, limit)

		if err != nil {
			return err
		}

		if format == FormatJSON {
			r.JSON(200, categories)
		} else {
			r.XML(200, models.Categories{Categories: categories})
		}

		return nil

	}).ServeHTTP(w, req)
}

func GetJWTQuestions(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render) {

	jwtHandlerFactory(func(w http.ResponseWriter, req *http.Request, token *jwt.Token) error {

		claims := utils.MapHelper{Map: token.Claims}
		format := claims.String(ParameterFormat, FormatXML)
		limit := claims.Int(ParameterLimit, DefaultLimit)
		offset := claims.Int(ParameterOffset, DefaultOffset)
		category := claims.String(ParameterCategory, "")

		if len(category) == 0 {
			http.Error(w, fmt.Errorf("Parameter '%s' is missing.", ParameterCategory).Error(), HttpStatusValidationError)
			return nil
		}

		questions, err := dm.QuestionStore.GetQuestions(req, limit,offset,category)

		if err != nil {
			return err
		}

		if format == FormatJSON {
			r.JSON(200, questions)
		} else {
			r.XML(200, models.Questions{Questions: questions})
		}

		return nil

	}).ServeHTTP(w, req)
}
