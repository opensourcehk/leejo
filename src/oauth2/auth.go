package oauth2

import (
	"github.com/RangelReale/osin"
	"log"
	"net/http"
	"upper.io/db"
)

func Bind(authPath string, sessPtr *db.Database) {

	// OAuth2 endpoints handler
	oauth2 := osin.NewServer(osin.NewServerConfig(), &AuthStorage{
		Db: *sessPtr,
	})

	// handle OAuth2 endpoints
	http.HandleFunc(authPath+"/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := oauth2.NewResponse()
		if ar := oauth2.HandleAuthorizeRequest(resp, r); ar != nil {
			// HANDLE LOGIN PAGE HERE
			// TODO:
			// 1. check if there is login data. If not, show form
			// 2. if there is login data from form, search for matched user
			// 3. if user exists, get user data and assign to ar.UserData, ar.Authorized = true
			ar.Authorized = true
			oauth2.FinishAuthorizeRequest(resp, r, ar)
		}
		log.Printf("OAuth2 Authorize Response: %#v", resp)
		osin.OutputJSON(resp, w, r)
	})
	http.HandleFunc(authPath+"/token", func(w http.ResponseWriter, r *http.Request) {
		resp := oauth2.NewResponse()
		if ar := oauth2.HandleAccessRequest(resp, r); ar != nil {
			log.Printf("Access successful")
			// TODO:
			// 1. check the auth code, client id and secret (or is it checked already?)
			// 2. if checking pass, generate and return token (or is it handled already?)
			ar.Authorized = true
			oauth2.FinishAccessRequest(resp, r, ar)
		} else {
			log.Printf("Internal Error: %s", resp.InternalError.Error())
		}
		log.Printf("OAuth2 Token Response: %#v", resp)
		osin.OutputJSON(resp, w, r)
	})

}
