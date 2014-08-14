package main

import (
	"log"
	"os"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Integration Tests Failed\n")
			os.Exit(1)
		}
	}()

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
