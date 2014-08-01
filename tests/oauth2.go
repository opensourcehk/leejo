package main

import (
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"log"
	"net/http"
	"net/url"
)

func getCode() (code string, err error) {
	ch := make(chan interface{})

	// temporary serve the redirect endpoint
	go TempServe(":8000", "/redirect/", ch,
		func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			log.Printf("Called handler function: %#v\n", q)

			// test if there is error
			if _, ok := q["error"]; ok {
				errStr := q.Get("error")
				fmt.Fprintf(w, "Error: Failed to login: %s\n", errStr)
			} else if _, ok := q["code"]; !ok {
				fmt.Fprintf(w, "Unknown error. Failed to obtain code\n")
			} else if _, ok := q["state"]; !ok {
				fmt.Fprintf(w, "Unknown error. Failed to obtain state\n")
			} else {
				fmt.Fprintf(w, "Login Success! Please go back to console to see test result\n")
			}

			ch <- q
			return
		})

	log.Printf("Test Auth\n")
	open.Start("http://localhost:8080/oauth2/authorize?" +
		"response_type=code&client_id=testing&scope=usersInfo")

	// wait for reply
	log.Printf("wait for result finish\n")
	res := <-ch
	result := res.(url.Values)

	// test if there is error
	if _, ok := result["error"]; ok {
		errStr := result.Get("error")
		log.Fatalf("Failed to login. Error: %s\n", errStr)
	} else if _, ok := result["code"]; !ok {
		log.Fatalf("Unknown error. Failed to obtain code\n")
	} else if _, ok := result["state"]; !ok {
		log.Fatalf("Unknown error. Failed to obtain state\n")
	}

	code = result.Get("code")
	state := result.Get("state")
	log.Printf("result: code=\"%s\" state=\"%s\"\n", code, state)

	return
}

func testAuth() (token string, err error) {

	// 1. use default browser to access authorize endpoint
	// 2. emulate redirect endpoint and retrieve code
	code, err := getCode()
	if err != nil {
		return
	}

	// access token endpoint
	// and retrieve token
	open.Start("http://localhost:8080/oauth2/token?" +
		"grant_type=authorization_code&code=" + code +
		"&client_id=testing&client_secret=testing" +
		"&redirect_uri=http://localhost:8000/redirect/")
	return
}
