package controllers

import (
	"github.com/martini-contrib/render"
)

func GetIndex(r render.Render) {
	GetIndexWithTemplateMap(r, nil)
}

func GetIndexWithTemplateMap(r render.Render, templateMap map[string]interface{}) {
	r.HTML(200, "index", templateMap)
}
