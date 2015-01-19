package main

import (
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/gourd/service"
	"github.com/gourd/session"
	"github.com/opensourcehk/leejo/lib/data"
	"github.com/opensourcehk/leejo/lib/rest"
	"log"
	"net/http"
	"net/url"
)

// basic login and authorization handler
type AuthHandler struct {
	UserHandler    rest.Handler
	SessionHandler session.Handler
}

func (h *AuthHandler) HandleLogin(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) (err error) {

	// parse POST input
	r.ParseForm()
	if r.Method == "POST" {
		// r.Form.Get("login") == "test" && r.Form.Get("password") == "test" {

		// get login information from form
		loginName := r.Form.Get("login")
		loginPass := r.Form.Get("password")
		if loginName == "" || loginPass == "" {
			err = fmt.Errorf("Empty Username or Password")
			return
		}

		// obtain session
		var sess session.Session
		sess, err = h.SessionHandler.Session(r)
		if err != nil {
			log.Printf("Session error: %s", err.Error())
			err = fmt.Errorf("Unknown server error")
			return
		}

		// obtain user service
		us := h.UserHandler.Service(sess)

		// get user from database
		var users []data.User
		c := service.NewConds().Add("username", loginName)
		err = us.Search(c, &users)
		if err != nil {
			log.Printf("Error searching user with Service: %s", err.Error())
			err = fmt.Errorf("Internal Server Error")
			return
		}

		// if user does not exists
		if len(users) == 0 {
			log.Printf("Unknown user \"%s\" attempt to login", loginName)
			err = fmt.Errorf("Username or Password incorrect")
			return
		}

		// if password does not match
		// TODO: use hash in password
		user := users[0]
		if user.PasswordHash != loginPass {
			log.Printf("Attempt to login \"%s\" with incorrect password", loginName)
			err = fmt.Errorf("Username or Password incorrect")
		}

		// return pointer of user object, allow it to be re-cast
		ar.UserData = &user
		return
	}

	// no POST input or incorrect login, show form
	err = fmt.Errorf("No login information")
	w.Write([]byte("<html><body>"))
	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"%s?response_type=%s&client_id=%s&state=%s&scope=%s&redirect_uri=%s\" method=\"POST\">",
		r.URL.Path,
		ar.Type,
		ar.Client.GetId(),
		ar.State,
		ar.Scope,
		url.QueryEscape(ar.RedirectUri))))
	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))
	w.Write([]byte("</form>"))
	w.Write([]byte("</body></html>"))
	return
}
