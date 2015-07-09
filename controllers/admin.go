package controllers

import (
	"fmt"
	"github.com/martini-contrib/render"
	"net/http"
	"services"
)

const (
	CategoryLimit         = 10
	TemplateKeyCategories = "Categories"
)

func Admin(q *services.QuizzicalService, w http.ResponseWriter, req *http.Request, r render.Render) {

	categories, err := q.CategoryStore.GetAll(req, CategoryLimit)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	} else {

		templateMap := make(map[string]interface{})
		templateMap[TemplateKeyCategories] = categories

		r.HTML(200, "admin", templateMap)
	}
}
