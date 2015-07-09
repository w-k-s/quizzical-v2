package controllers

import (
	"github.com/martini-contrib/render"
	"models"
	"net/http"
	"services"
)

func PostCategory(q *services.QuizzicalService, postedCategory models.Category, w http.ResponseWriter, req *http.Request, r render.Render) {

	err := q.CategoryStore.Save(req, &postedCategory)

	templateMap := make(map[string]interface{})

	if err != nil {
		templateMap[TemplateKeyCategoryError] = err.Error()
	} else {
		templateMap[TemplateKeyCategory] = postedCategory
	}

	GetAdminWithTemplateMap(q, w, req, r, templateMap)

}
