package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/gorilla/pat"
	"io/ioutil"
	"leejo/oauth2"
	session "github.com/gourd/session/oauth2"
	"log"
	"net/http"
	"os"
	"upper.io/db"
	"upper.io/db/postgresql"
)

type config struct {
	Db dbconfig `json:"db"`
}

type dbconfig struct {
	Host string `json:"host"`
	Name string `json:"name"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func main() {

	// parse flags
	confFn := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	// read the config file to conf
	confFile, err := ioutil.ReadFile(*confFn)
	if err != nil {
		fmt.Printf("Failed opening config file \"%s\": %v\n", *confFn, err)
		os.Exit(1)
	}
	conf := config{}
	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		fmt.Printf("Failed parsing config file \"%s\": %v\n", *confFn, err)
		os.Exit(1)
	}

	// logs to stdout
	log.SetOutput(os.Stdout)

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

	// oauth2 related config
	oConf := osin.NewServerConfig()
	oConf.AllowGetAccessRequest = true
	oConf.AllowClientSecretInParams = true
	oConf.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
	}
	oConf.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{
		osin.CODE,
		osin.TOKEN,
	}

	// oauth2 endpoints handler
	oStore := &oauth2.AuthStorage{
		Db: dbs,
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

	// map API to pat router
	RestOnPat(uh, sh, r)
	RestOnPat(ush, sh, r)
	RestOnPat(uih, sh, r)

	// handle OAuth2 endpoints
	oauth2.BindOsin("/oauth2", osin.NewServer(oConf, oStore), lh)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

}
