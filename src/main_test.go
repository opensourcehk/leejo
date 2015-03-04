package main

import (
	"net/http/httptest"
	"testing"
	"upper.io/db"
	"upper.io/db/postgresql"
)

func testMain(t *testing.T) {

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

	// start the server
	ts := httptest.NewServer(createHandler(dbs))
	defer ts.Close()

}
