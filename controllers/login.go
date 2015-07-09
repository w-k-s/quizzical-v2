package controllers

import (
	_ "appengine"
	"auth"
	"fmt"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"models"
	"net/http"
)

func Login(session sessions.Session, postedUser models.User, r render.Render, w http.ResponseWriter, req *http.Request) {

	if !auth.Authenticate(postedUser.Username, postedUser.Password) {

		r.Redirect(sessionauth.RedirectUrl)
		return
	}

	//ugly hack
	postedUser.Id = auth.MasterUserId

	err := sessionauth.AuthenticateSession(session, &postedUser)
	if err != nil {

		fmt.Fprintf(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := req.URL.Query()
	redirect := params.Get(sessionauth.RedirectParam)

	r.Redirect(redirect)
}

func Logout(session sessions.Session, user sessionauth.User, r render.Render) {
	sessionauth.Logout(session, user)
	r.Redirect("/")
}
