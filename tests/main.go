package main

import (
	"log"
)

func main() {

	// test authentication service
	token, err := testAuth()
	if err != nil {
		panic(err)
	}

	// test APIs
	err = testUser(token)
	if err != nil {
		panic(err)
	}

	log.Printf("Integration Tests Pass\n")

}
