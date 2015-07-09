package controllers

import (
	"github.com/martini-contrib/render"
	"models"
	"net/http"
	"services"
)

func PostQuestion(q *services.QuizzicalService, postedQuestion models.Question, w http.ResponseWriter, req *http.Request, r render.Render) {

	err := q.QuestionStore.Save(req, &postedQuestion)

	templateMap := make(map[string]interface{})

	if err != nil {
		templateMap[TemplateKeyQuestionError] = err.Error()
	} else {
		templateMap[TemplateKeyQuestion] = postedQuestion
	}

	GetAdminWithTemplateMap(q, w, req, r, templateMap)

}
