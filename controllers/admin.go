package controllers

import (
	"datastore"
	"fmt"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net/http"
	"utils"
)

func GetAdmin(dm *datastore.Manager, session sessions.Session, w http.ResponseWriter, req *http.Request, r render.Render) {

	categories, err := dm.CategoryStore.GetAll(req, -1)

	if err != nil {
		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
	} else {

		templateMap := make(map[string]interface{})

		templateMap[TemplateKeyCategories] = categories
		templateMap[TemplateKeyCategoryError] = utils.PopFlash(session, TemplateKeyCategoryError)
		templateMap[TemplateKeyCategory] = utils.PopFlash(session, TemplateKeyCategory)
		templateMap[TemplateKeyQuestion] = utils.PopFlash(session, TemplateKeyQuestion)

		r.HTML(200, "admin", templateMap)
	}
}
