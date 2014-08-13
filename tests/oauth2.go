package main

import (
	"encoding/json"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// serve redirection endpoints for test
func serveRedirectEnd(addr string, basePath string) (ch chan chanResp) {

	ch = make(chan chanResp)

	// temporary serve the redirect endpoint
	go func() {

		// handler for response_type=authentication_code test
		http.HandleFunc(basePath+"authcode/", func(w http.ResponseWriter, r *http.Request) {
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
				fmt.Fprintf(w, "Login Success! Please go back to console "+
					"to see test result\n")
			}

			ch <- chanResp{
				Data: q,
			}
			return
		})

		// handler for response_type=token test
		http.HandleFunc(basePath+"token/", func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			log.Printf("Called handler function 2: %#v\n", q)

			// test if there is error
			fmt.Fprintf(w, "<!DOCTYPE html>"+
				"<html><head><script>"+
				"  if (window.location.hash) {"+
				"    console.log(window.location.hash);"+
				"    var q = window.location.hash.substr(1);"+
				"    window.location.href = '"+basePath+"token/result/?' + q;"+
				"  }"+
				"</script></head>"+
				"<body><b>Please enable javascript and try again</b></body></html>\n")

			return
		})

		// handler for response_type=authentication_code test
		http.HandleFunc(basePath+"token/result/", func(w http.ResponseWriter, r *http.Request) {
			q := r.URL.Query()
			log.Printf("Called handler function 3: %#v\n", q)

			// test if there is error
			if _, ok := q["error"]; ok {
				errStr := q.Get("error")
				fmt.Fprintf(w, "Error: Failed to login: %s\n", errStr)
			} else if _, ok := q["access_token"]; !ok {
				fmt.Fprintf(w, "Unknown error. Failed to obtain access token\n")
			} else {
				fmt.Fprintf(w, "Login Success! Please go back to console "+
					"to see test result\n")
			}

			ch <- chanResp{
				Data: q,
			}
			return
		})

		log.Fatal(http.ListenAndServe(addr, nil))

	}()

	return
}

// obtain authorization code
func getAuthCodeResp(ch chan chanResp) (code string, err error) {

	log.Printf("Test Auth\n")

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", "testing")
	params.Set("scope", "user,user_skills,user_interests")
	params.Set("redirect_uri", "http://localhost:8000/redirect/authcode/")
	open.Start("http://localhost:8080/oauth2/authorize?" + params.Encode())

	// wait for reply
	log.Printf("wait for result finish\n")
	res := <-ch
	if res.Err != nil {
		log.Fatalf("Failed: Error: %s\n", res.Err.Error())
	}

	// cast result data
	result := res.Data.(url.Values)

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

// obtain authorization code
func getTokenResp(ch chan chanResp) (tResp tokenResp, err error) {

	log.Printf("Test Auth (response_type=token)\n")

	params := url.Values{}
	params.Set("response_type", "token")
	params.Set("client_id", "testing")
	params.Set("scope", "user,user_skills,user_interests")
	params.Set("redirect_uri", "http://localhost:8000/redirect/token/")
	open.Start("http://localhost:8080/oauth2/authorize?" + params.Encode())

	// wait for reply
	log.Printf("wait for result finish\n")
	resp := <-ch
	respVal := resp.Data.(url.Values)

	// decode the json response
	tResp.ErrMsg = respVal.Get("error")
	tResp.ErrDesc = respVal.Get("error_description")
	tResp.AccessToken = respVal.Get("access_token")
	expiresStr := respVal.Get("expires_in")
	var i int64
	if expiresStr != "" {
		i, err = strconv.ParseInt(expiresStr, 10, 64)
	}
	tResp.ExpiresIn = i
	tResp.Scope = respVal.Get("scope")
	tResp.TokenType = respVal.Get("token_type")
	log.Printf("getTokenResp %#v\n", tResp)

	// if the request do not result in error
	if tResp.HasError() {
		log.Fatalf("Failed to login. Error: %s\n", tResp.Error())
		err = &tResp
		return
	}
	return
}

type chanResp struct {
	Data interface{}
	Err  error
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
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("client_id", "testing")
	params.Set("client_secret", "testing")
	params.Set("code", code)
	params.Set("redirect_uri", "http://localhost:8000/redirect/authcode/")
	req, err := http.NewRequest("GET",
		"http://localhost:8080/oauth2/token?"+params.Encode(), nil)

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
	req, err := http.NewRequest("GET", "http://localhost:8080/oauth2/token?"+
		"grant_type=refresh_token"+
		"&client_id=testing&client_secret=testing"+
		"&refresh_token="+refreshT+
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

	// serve up redirection endpoint(s)
	ch := serveRedirectEnd(":8000", "/redirect/")

	// 1. use default browser to access authorize endpoint
	// 2. emulate redirect endpoint and retrieve code
	code, err := getAuthCodeResp(ch)
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

	// test getting access token with response_type=token
	tResp, err = getTokenResp(ch)
	if err != nil {
		return
	}

	// return only the access token for other tests
	accessT = tResp.AccessToken
	return
}
