package main

import (
	"github.com/gorilla/pat"
	session "github.com/gourd/session/oauth2"
	"leejo/oauth2"
	"leejo/rest"
	"log"
	"net/http"
	"os"
	"upper.io/db"
	"upper.io/db/postgresql"
)

// set initial envirnment
func init() {

	// logs to stdout
	log.SetOutput(os.Stdout)

}

func main() {

	// read config files
	conf := getConfig()

	// connect to database
	var dbsettings = db.Settings{
		Host:     conf.Db.Host,
		Database: conf.Db.Name,
		User:     conf.Db.User,
		Password: conf.Db.Pass,
	}
	dbs, err := db.Open(postgresql.Adapter, dbsettings)
	if err != nil {
		panic(err)
	}
	defer dbs.Close()

	// Users REST API helper
	uh := &UserRest{
		Db:       dbs,
		basePath: "/api.v1/users",
		subPath:  "{id:[0-9]+}",
	}

	// UserSkills REST API helper
	ush := &UserSkillRest{
		Db:       dbs,
		basePath: "/api.v1/userSkills/{user_id:[0-9]+}",
		subPath:  "{id:[0-9]+}",
	}

	// UserInterests REST API helper
	uih := &UserInterestRest{
		Db:       dbs,
		basePath: "/api.v1/userInterests/{user_id:[0-9]+}",
		subPath:  "{id:[0-9]+}",
	}

	// oauth2 endpoints handler
	oStore := &oauth2.AuthStorage{
		// provides services related to oauth2
		P: &oauth2Provider{
			Db: dbs,
		},
	}

	// define session handler
	// that works with a osin server
	sh := &session.OsinHandler{
		Storage: oStore,
	}

	// handle login
	lh := &AuthHandler{
		SessionHandler: sh,
		UserHandler:    uh,
	}

	// gorilla pat for routing
	r := pat.New()

	// protocol
	p := &Protocol{}

	// map API to pat router
	rest.Pat(uh, sh, p, r)
	rest.Pat(ush, sh, p, r)
	rest.Pat(uih, sh, p, r)

	// handle OAuth2 endpoints
	oauth2.Pat("/oauth2", oStore, sh, lh, r)

	// start the server
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

}
