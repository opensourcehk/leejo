package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/RangelReale/osin"
	"github.com/gorilla/pat"
	"io/ioutil"
	"leejo/oauth2"
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

	// connect to database
	var dbsettings = db.Settings{
		Host:     conf.Db.Host,
		Database: conf.Db.Name,
		User:     conf.Db.User,
		Password: conf.Db.Pass,
	}
	sess, err := db.Open(postgresql.Adapter, dbsettings)
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	// oauth2 related config
	osinConf := osin.NewServerConfig()
	osinConf.AllowGetAccessRequest = true
	osinConf.AllowClientSecretInParams = true
	osinConf.AllowedAccessTypes = osin.AllowedAccessType{
		osin.AUTHORIZATION_CODE,
		osin.REFRESH_TOKEN,
	}
	osinConf.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{
		osin.CODE,
		osin.TOKEN,
	}

	// OAuth2 endpoints handler
	osinServer := osin.NewServer(osinConf, &oauth2.AuthStorage{
		Db: sess,
	})

	// gorilla pat for routing
	r := pat.New()

	// Users related API
	bindUser("/api.v1/users", osinServer, &sess, r)

	// UserSkills related API
	bindUserSkills("/api.v1/userSkills/{user_id:[0-9]+}", osinServer, &sess, r)

	// UserInterests related API
	bindUserInterests("/api.v1/userInterests/{user_id:[0-9]+}", osinServer, &sess, r)

	// handle OAuth2 endpoints
	oauth2.Bind("/oauth2", osinServer)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

}
