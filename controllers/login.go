package controllers

import(
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/sessions"
	"models"
	"auth"
	"net/http"
	"fmt"
)

func Login(session sessions.Session, postedUser models.User, r render.Render, w http.ResponseWriter, req *http.Request){

	if !auth.Authenticate(postedUser.Username, postedUser.Password){
		r.Redirect(sessionauth.RedirectUrl)
		return
	}

	err := sessionauth.AuthenticateSession(session,&postedUser)
	
	if err != nil {
		fmt.Fprintf(w,err.Error(),http.StatusInternalServerError)
	}

	params := req.URL.Query()
	redirect := params.Get(sessionauth.RedirectParam)
	r.Redirect(redirect)
}