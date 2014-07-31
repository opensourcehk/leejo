package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {

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

	// test authentication service
	_, err := testAuth()
	if err != nil {
		panic(err)
	}

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

	code := result.Get("code")
	state := result.Get("state")
	log.Printf("result: code=\"%s\" state=\"%s\"\n", code, state)

	// test APIs
	err = testUser()
	if err != nil {
		panic(err)
	}

	log.Printf("Integration Tests Pass\n")

}
