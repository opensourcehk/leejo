package main

import (
	"log"
)

func main() {

	// test authentication service
	_, err := testAuth()
	if err != nil {
		panic(err)
	}

	// test APIs
	err = testUser()
	if err != nil {
		panic(err)
	}

	log.Printf("Integration Tests Pass\n")

}
