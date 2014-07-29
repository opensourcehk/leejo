package main

import (
	"github.com/RangelReale/osin"
	"net/http"
	"upper.io/db"
)

func bindAuth(authPath string, sessPtr *db.Database) {

	// OAuth2 endpoints handler
	oauth2 := osin.NewServer(osin.NewServerConfig(), &AuthStorage{
		Db: *sessPtr,
	})

	// handle OAuth2 endpoints
	http.HandleFunc(authPath+"/authorize", func(w http.ResponseWriter, r *http.Request) {
		resp := oauth2.NewResponse()
		if ar := oauth2.HandleAuthorizeRequest(resp, r); ar != nil {
			// HANDLE LOGIN PAGE HERE
			ar.Authorized = true
			oauth2.FinishAuthorizeRequest(resp, r, ar)
		}
		osin.OutputJSON(resp, w, r)
	})
	http.HandleFunc(authPath+"/token", func(w http.ResponseWriter, r *http.Request) {
		resp := oauth2.NewResponse()
		if ar := oauth2.HandleAccessRequest(resp, r); ar != nil {
			ar.Authorized = true
			oauth2.FinishAccessRequest(resp, r, ar)
		}
		osin.OutputJSON(resp, w, r)
	})

}
