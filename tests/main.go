package main

import (
	"fmt"
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

	fmt.Printf("Integration Tests Pass\n")

}
