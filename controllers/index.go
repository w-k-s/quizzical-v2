package controllers

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"utils"
)

func GetIndex(session sessions.Session, r render.Render) {

	templateMap := make(map[string]interface{})

	templateMap[TemplateKeyAuthenticationFailed] = utils.PopFlash(session, TemplateKeyAuthenticationFailed)

	r.HTML(200, "index", templateMap)
}
