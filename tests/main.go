package main

import (
	"fmt"
)

func main() {
	err := testUser()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Integration Tests Pass\n")
}
