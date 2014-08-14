package main

import (
	"fmt"
	"github.com/RangelReale/osin"
	"net/http"
	"net/url"
)

// basic login and authorization handler
type AuthHandler struct {
}

func (h *AuthHandler) HandleLogin(ar *osin.AuthorizeRequest, w http.ResponseWriter, r *http.Request) (err error) {

	// parse POST input
	r.ParseForm()
	if r.Method == "POST" && r.Form.Get("login") == "test" && r.Form.Get("password") == "test" {
		// TODO: really handle the login with database data
		somedata := "some data"
		ar.UserData = &somedata
		return
	}

	// no POST input or incorrect login, show form
	err = fmt.Errorf("No login information")
	w.Write([]byte("<html><body>"))
	w.Write([]byte(fmt.Sprintf("LOGIN %s (use test/test)<br/>", ar.Client.GetId())))
	w.Write([]byte(fmt.Sprintf("<form action=\"%s?response_type=%s&client_id=%s&state=%s&redirect_uri=%s\" method=\"POST\">",
		r.URL.Path,
		ar.Type,
		ar.Client.GetId(),
		ar.State,
		url.QueryEscape(ar.RedirectUri))))
	w.Write([]byte("Login: <input type=\"text\" name=\"login\" /><br/>"))
	w.Write([]byte("Password: <input type=\"password\" name=\"password\" /><br/>"))
	w.Write([]byte("<input type=\"submit\"/>"))
	w.Write([]byte("</form>"))
	w.Write([]byte("</body></html>"))
	return
}
