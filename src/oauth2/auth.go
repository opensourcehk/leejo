package oauth2

import (
	"github.com/RangelReale/osin"
	"log"
	"net/http"
)

// interface to handle login and authorization
type AuthHandler interface {
	HandleLogin(*osin.AuthorizeRequest, http.ResponseWriter, *http.Request) error
}

// bind the endpoints to http server
func BindOsin(authPath string, osinServer *osin.Server, lh AuthHandler) {

	// handle OAuth2 endpoints
	http.HandleFunc(authPath+"/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := osinServer.NewResponse()
		if ar := osinServer.HandleAuthorizeRequest(resp, r); ar != nil {
			// handle login page
			if err := lh.HandleLogin(ar, w, r); err != nil {
				return
			}
			log.Printf("OAuth2 Authorize Request: User obtained: %#v", ar.UserData)
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
			// TODO: handle authorization
			// check if the user has the permission to grant the scope
			log.Printf("Access successful")
			ar.Authorized = true
			osinServer.FinishAccessRequest(resp, r, ar)
		} else if resp.InternalError != nil {
			log.Printf("Internal Error: %s", resp.InternalError.Error())
		}
		log.Printf("OAuth2 Token Response: %#v", resp)
		osin.OutputJSON(resp, w, r)
	})

}
