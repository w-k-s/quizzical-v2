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
	"errors"
	JWT "bitbucket.org/waqqas-abdulkareem/asfour_toolkit/handlers/jwt"
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

//---------------------------------- API v2 -----------------------------------//

type ResponseFormatter func(r * http.Request,w http.ResponseWriter, response interface{}, err error)

type QuizzicalApi struct{
	Consumer *jwt.Consumer
	DB *datastore.Manager
	ResponseFormatter ResponseFormatter
}

func (api * QuizzicalApi) Categories(r * http.Request, w http.ResponseWriter){

	var resp interface{}

	token,err := api.Consumer.ProcessTokenFromRequestParameter(r,"token",[]byte(auth.JWTSecret))

	if err == nil {

		limit := token.Int32("limit",DefaultLimit)
		categories,err := api.DB.CategoryStore.GetAll(r,int(limit))

		if err == nil{
			resp = categories;
		}
	}

	api.ResponseFormatter(r,w,resp,err)
}

func (api * QuizzicalApi) PostCategory(r * http.Request, w http.ResponseWriter){

	var resp interface{}

	token, err := api.Consumer.ProcessTokenFromRequestParameter(r,"token",[]byte(auth.JWTSecret))

	if err != nil {
		api.ResponseFormatter(r,w,nil,err)
		return
	}

	name := token.String("name","")

	if len(name) == 0 {
		api.ResponseFormatter(r,w,nil,errors.New("Name can not be empty"))
		return
	}

	category := models.Category{Name: name}
	err = api.DB.CategoryStore.Save(r,&category)
	resp = category;

	api.ResponseFormatter(r,w,resp,err)
}

func jwtHandlerFactory(handler func(http.ResponseWriter, *http.Request, *jwt.Token) error) *JWT.Handler {

	return JWT.NewHandler(
		"token",
		[]byte(auth.JWTSecret),
		handler,
	)
}

//--------------------------------------------  API V1 -----------------------------------------------//

func GetJWTCategories(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render) {

	jwtHandlerFactory(func(w http.ResponseWriter, req *http.Request, token *jwt.Token) error {

		claims := utils.MapHelper{Map: token.Claims}
		format := claims.String(ParameterFormat, FormatXML)

		limit := DefaultLimit
		if _limit,ok := token.Claims[ParameterLimit]; ok{
			limit = int(_limit.(float64))
		}

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
		category := claims.String(ParameterCategory, "")

		if len(category) == 0 {
			http.Error(w, fmt.Errorf("Parameter '%s' is missing.", ParameterCategory).Error(), HttpStatusValidationError)
			return nil
		}

		limit := DefaultLimit
		if _limit,ok := token.Claims[ParameterLimit]; ok{
			limit = int(_limit.(float64))
		}

		/*offset := DefaultOffset
		if _offset,ok := token.Claims[ParameterOffset];ok{
			offset = int(_offset.(float64))
		}*/

		questions, err := dm.QuestionStore.Random(req, category, limit)

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
