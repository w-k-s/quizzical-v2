package controllers

import (
	"github.com/martini-contrib/render"
)

func Index(r render.Render) {
	r.HTML(200, "index", "")
}
