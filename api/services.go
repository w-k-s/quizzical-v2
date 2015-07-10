package api

import (
	"datastore"
	"fmt"
	"github.com/martini-contrib/render"
	"models"
	"net/http"
	"strconv"
	"utils"
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

func GetCategories(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render) {

	format := utils.FormValue(req, ParameterFormat, FormatXML)
	limit, _ := strconv.Atoi(utils.FormValue(req, ParameterLimit, DefaultLimitString))

	categories, err := dm.CategoryStore.GetAll(req, limit)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	} else if format == FormatJSON {
		r.JSON(200, categories)
	} else {
		r.XML(200, models.Categories{Categories: categories})
	}
}

func GetQuestions(dm *datastore.Manager, w http.ResponseWriter, req *http.Request, r render.Render) {

	format := utils.FormValue(req, ParameterFormat, FormatXML)
	limit, _ := strconv.Atoi(utils.FormValue(req, ParameterLimit, DefaultLimitString))
	category := req.FormValue(ParameterCategory)

	if len(category) == 0 {
		fmt.Fprintf(w, fmt.Errorf("Parameter '%s' is missing.", ParameterCategory).Error(), HttpStatusValidationError)
		return
	}

	questions, err := dm.QuestionStore.Random(req, category, limit)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	} else if format == FormatJSON {
		r.JSON(200, questions)
	} else {
		r.XML(200, models.Questions{Questions: questions, Category: category})
	}
}
