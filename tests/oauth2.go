package main

import (
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"log"
	"net/http"
	"net/url"
)

// obtain authorization code
func getAuthCode() (code string, err error) {
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

type tokenResp struct {
	ErrMsg       string `json:"error,omitempty"`
	ErrDesc      string `json:"error_description,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	Scope        string `json:"scope,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
}

func (r *tokenResp) HasError() bool {
	return r.Error() != ""
}

func (r *tokenResp) Error() string {
	return r.ErrMsg
}


// obtain token with authorization code
func getToken(code string) (tResp tokenResp, err error) {

	// generate token request
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/oauth2/token?" +
		"grant_type=authorization_code&code=" + code +
		"&client_id=testing&client_secret=testing" +
		"&redirect_uri=http://localhost:8000/redirect/", nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	// decode the json response
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&tResp)
	log.Printf("getToken %#v\n", tResp)

	// if the request do not result in error
	if tResp.HasError() {
		log.Fatalf("Failed to login. Error: %s\n", tResp.Error())
		err = &tResp
		return
	}
	return
}

// test to refresh token
func refreshToken(refreshT string) (tResp tokenResp, err error) {
	// generate token request
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/oauth2/token?" +
		"grant_type=refresh_token" +
		"&client_id=testing&client_secret=testing" +
		"&refresh_token=" + refreshT +
		"&redirect_uri=http://localhost:8000/redirect/", nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	// decode the json response
	dec := json.NewDecoder(resp.Body)
	dec.Decode(&tResp)
	log.Printf("refreshToken %#v\n", tResp)

	// if the request do not result in error
	if tResp.HasError() {
		log.Fatalf("Failed to login. Error: %s\n", tResp.Error())
		err = &tResp
		return
	}
	return
}

// main test on authentication and authorization routine
func testAuth() (accessT string, err error) {

	// 1. use default browser to access authorize endpoint
	// 2. emulate redirect endpoint and retrieve code
	code, err := getAuthCode()
	if err != nil {
		return
	}

	// access token endpoint
	// and retrieve token
	tResp, err := getToken(code)
	if err != nil {
		return
	}

	// access token endpoint
	// and retrieve token
	tResp, err = refreshToken(tResp.RefreshToken)
	if err != nil {
		return
	}

	// return only the access token for other tests
	accessT = tResp.AccessToken
	return
}
