package main

import (
	"log"
	"net/http"
)

func testAuth() (auth interface{}, err error) {
	log.Printf("Test Auth\n")
	_, err =
		http.Get("http://localhost:8080/oauth2/authorize?response_type=code&client_id=testing")
	return
}
