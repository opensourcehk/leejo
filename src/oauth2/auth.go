package oauth2

import (
	"github.com/RangelReale/osin"
	"log"
	"net/http"
)

func Bind(authPath string, osinServer *osin.Server) {

	// handle OAuth2 endpoints
	http.HandleFunc(authPath+"/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := osinServer.NewResponse()
		if ar := osinServer.HandleAuthorizeRequest(resp, r); ar != nil {
			// HANDLE LOGIN PAGE HERE
			// TODO:
			// 1. check if there is login data. If not, show form
			// 2. if there is login data from form, search for matched user
			// 3. if user exists, get user data and assign to ar.UserData, ar.Authorized = true
			ar.Authorized = true
			osinServer.FinishAuthorizeRequest(resp, r, ar)
		}
		if resp.InternalError != nil {
			log.Printf("Internal Error: %s", resp.InternalError.Error())
		}
		log.Printf("OAuth2 Authorize Response: %#v", resp)
		osin.OutputJSON(resp, w, r)
	})
	http.HandleFunc(authPath+"/token", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Access token endpoint")
		resp := osinServer.NewResponse()

		// ugly hack to accept GET request in token endpoint
		// should add this to osin
		if len(r.Form) == 0 {
			r.Form = r.URL.Query()
		}

		if ar := osinServer.HandleAccessRequest(resp, r); ar != nil {
			log.Printf("Access successful")
			// TODO:
			// 1. check the auth code, client id and secret (or is it checked already?)
			// 2. if checking pass, generate and return token (or is it handled already?)
			ar.Authorized = true
			osinServer.FinishAccessRequest(resp, r, ar)
		} else if resp.InternalError != nil {
			log.Printf("Internal Error: %s", resp.InternalError.Error())
		}
		log.Printf("OAuth2 Token Response: %#v", resp)
		osin.OutputJSON(resp, w, r)
	})

}
