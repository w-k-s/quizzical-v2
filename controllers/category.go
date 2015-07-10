package controllers

import (
	"datastore"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"models"
	"net/http"
	"utils"
)

func PostCategory(session sessions.Session, dm *datastore.Manager, postedCategory models.Category, w http.ResponseWriter, req *http.Request, r render.Render) {

	err := dm.CategoryStore.Save(req, &postedCategory)

	if err != nil {
		utils.PushFlash(session, TemplateKeyCategoryError, err.Error())
	} else {
		utils.PushFlash(session, TemplateKeyCategory, &postedCategory)
	}

	r.Redirect("/admin")
	return
}
