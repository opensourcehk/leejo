package main

import (
	"github.com/skratchdot/open-golang/open"
	"log"
	//"net/http"
)

func testAuth() (auth interface{}, err error) {
	log.Printf("Test Auth\n")
	open.Start("http://localhost:8080/oauth2/authorize?response_type=code&client_id=testing")
	return
}
