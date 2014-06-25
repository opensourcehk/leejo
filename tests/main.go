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
	var user User

	p := napping.Params{}

	userToCreate := map[string]string{
		"username": "Tester", // TODO: use uuid
		"gender":   "F",
	}
	userToUpdate := map[string]string{
		"username": "Tester Updated",
		"gender":   "M",
	}

	// -- Create Test --
	// create a user
	resp, err = napping.Post("http://localhost:8080/api.v1/user",
		&userToCreate, &result, nil)

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

	// retrieve the user just created
	resp, err = napping.Get(
		fmt.Sprintf("http://localhost:8080/api.v1/user/%d",
			userId),
		&p, &result, nil)
	if err != nil {
		return
	}
	resultNum = len(result.Result)
	if resultNum != 1 {
		fmt.Printf("Raw: %s\n", resp.RawText())
		return fmt.Errorf("Bad response in retrieve user. "+
			"There are %d results (expecting 1)",
			resultNum)
	}

	// test: Test the retrieved user data
	user = result.Result[0]
	if user.Username != userToCreate["username"] {
		fmt.Printf("Raw: %s\n", resp.RawText())
		return fmt.Errorf("User create error. Username is %s (expected %s)",
			user.Username, userToCreate["username"])
	} else if user.Gender != userToCreate["gender"] {
		fmt.Printf("Raw: %s\n", resp.RawText())
		return fmt.Errorf("User create error. Gender is %s (expected %s)",
			user.Username, userToCreate["gender"])
	}

	// -- Update Test --
	// update the user just created
	resp, err = napping.Put(
		fmt.Sprintf("http://localhost:8080/api.v1/user/%d",
			userId),
		&userToUpdate, &result, nil)
	if err != nil {
		return
	}

	// retrieve the user just updated
	resp, err = napping.Get(
		fmt.Sprintf("http://localhost:8080/api.v1/user/%d",
			userId),
		&p, &result, nil)
	if err != nil {
		return
	}

	// test: Test the retrieved user data
	user = result.Result[0]
	if user.Username != userToUpdate["username"] {
		fmt.Printf("Raw: %s\n", resp.RawText())
		return fmt.Errorf("User create error. Username is %s (expected %s)",
			user.Username, userToUpdate["username"])
	} else if user.Gender != userToUpdate["gender"] {
		fmt.Printf("Raw: %s\n", resp.RawText())
		return fmt.Errorf("User create error. Gender is %s (expected %s)",
			user.Username, userToUpdate["gender"])
	}

	// -- Delete Test --
	// test: delete the user just created
	resp, err = napping.Delete(
		fmt.Sprintf("http://localhost:8080/api.v1/user/%d",
			userId),
		&result, nil)
	if err != nil {
		return
	}

	return
}

func main() {
	err := testUser()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Integration Tests Pass\n")
}
