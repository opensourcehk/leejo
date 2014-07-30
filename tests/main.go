package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	ch := make(chan interface{})

	// temporary serve the redirect endpoint
	go TempServe(":8000", "/redirect/", ch,
		func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Called handler function\n")
		ch <- nil
		fmt.Fprintf(w, "Hello world from my Go program!\n")
		return
	})

	// test authentication service
	_, err := testAuth()
	if err != nil {
		panic(err)
	}

	// wait for reply
	log.Printf("wait for result finish\n")
	result := <-ch
	log.Printf("result: %#v\n", result)

	// test APIs
	err = testUser()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Integration Tests Pass\n")

}
