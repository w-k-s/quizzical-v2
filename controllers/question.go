package controllers

import (
	"datastore"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"models"
	"net/http"
	"utils"
)

func PostQuestion(dm *datastore.Manager, session sessions.Session, postedQuestion models.Question, w http.ResponseWriter, req *http.Request, r render.Render) {

	err := dm.QuestionStore.Save(req, &postedQuestion)

	if err != nil {
		utils.PushFlash(session, TemplateKeyQuestionError, err.Error())
	} else {
		utils.PushFlash(session, TemplateKeyQuestion, &postedQuestion)
	}

	r.Redirect("/admin")

}
