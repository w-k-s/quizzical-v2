package controllers

import (
	"fmt"
	"github.com/martini-contrib/render"
	"net/http"
	"services"
)

func GetAdmin(q *services.QuizzicalService, w http.ResponseWriter, req *http.Request, r render.Render) {

	GetAdminWithTemplateMap(q, w, req, r, nil)
}

func GetAdminWithTemplateMap(q *services.QuizzicalService, w http.ResponseWriter, req *http.Request, r render.Render, templateMap map[string]interface{}) {

	if templateMap == nil {
		templateMap = make(map[string]interface{})
	}

	categories, err := q.CategoryStore.GetAll(req, -1)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	} else {

		templateMap[TemplateKeyCategories] = categories

		r.HTML(200, "admin", templateMap)
	}
}
