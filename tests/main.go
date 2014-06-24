package main

import (
	"fmt"
	"github.com/jmcvetta/napping"
)

type UserResp struct {
	Status string `json:"status"`
	Result []User `json:"result"`
}

func testUser() (err error) {

	var result UserResp
	var resp *napping.Response
	var resultNum int

	p := napping.Params{}

	// create a user
	resp, err = napping.Post("http://localhost:8080/api.v1/user",
		&map[string]string{
			"username": "Tester", // TODO: use uuid
			"gender":   "F",
		}, &result, nil)

	resultNum = len(result.Result)
	if resultNum != 1 {
		fmt.Printf("Raw: %s\n", resp.RawText())
		return fmt.Errorf("Bad response in create user. "+
			"There are %d results (expecting 1)",
			resultNum)
	}
	fmt.Printf("Success creating user\n")

	// retrieve the user
	resp, err = napping.Get("http://localhost:8080/api.v1/user/1",
		&p, &result, nil)
	if err != nil {
		return
	}
	if resp.Status() == 404 {
		err = fmt.Errorf("I don't know why")
		return
	}
	fmt.Printf("Integration test %d\n", resp.Status())
	return
}

func main() {
	err := testUser()
	if err != nil {
		panic(err)
	}
}
