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

	// test: has to be 1 row
	resultNum = len(result.Result)
	if resultNum != 1 {
		fmt.Printf("Raw: %s\n", resp.RawText())
		return fmt.Errorf("Bad response in create user. "+
			"There are %d results (expecting 1)",
			resultNum)
	}

	// test: the id should not be 0
	userId := result.Result[0].UserId
	if userId == 0 {
		fmt.Printf("Raw: %s\n", resp.RawText())
		return fmt.Errorf("Bad response in create user. " +
			"The returned user has a UserId = 0")
	}

	fmt.Printf("Success creating user\n")

	// retrieve the user just created
	resp, err = napping.Get(
		fmt.Sprintf("http://localhost:8080/api.v1/user/%d",
			userId),
		&p, &result, nil)
	if err != nil {
		return
	}

	// test: delete the user just created
	resp, err = napping.Delete(
		fmt.Sprintf("http://localhost:8080/api.v1/user/%d",
			userId),
		&result, nil)
	if err != nil {
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
