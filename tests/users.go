package main

import (
	"fmt"
	"github.com/jmcvetta/napping"
	"github.com/yookoala/restit"
)

type UserResp struct {
	Status string `json:"status"`
	Result []User `json:"result"`
}

func (r *UserResp) Count() int {
	return len(r.Result)
}

func (r *UserResp) NthExists(n int) (err error) {
	if n < 0 || n > r.Count() {
		err = fmt.Errorf("Nth item (%d) not exist. Length = %d",
			n, len(r.Result))
	}
	return
}

func (r *UserResp) NthValid(n int) (err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// test: the id should not be 0
	userId := r.Result[n].UserId
	if userId == 0 {
		return fmt.Errorf("The user has a UserId = 0")
	}

	return
}

func (r *UserResp) NthMatches(n int, comp *interface{}) (err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// check if the item match the payload
	user := r.Result[n]
	cptr := (*comp).(*map[string]string)
	c := *cptr
	if user.Username != c["username"] {
		err = fmt.Errorf("Username is \"%s\" (expected \"%s\")",
			user.Username, c["username"])
		return
	} else if user.Gender != c["gender"] {
		err = fmt.Errorf("Gender is \"%s\" (expected \"%s\")",
			user.Gender, c["gender"])
		return
	}

	return
}

func testUser() (err error) {

	var result UserResp
	var resp *napping.Response

	userToCreate := map[string]string{
		"username": "Tester", // TODO: use uuid
		"gender":   "F",
	}
	userToUpdate := map[string]string{
		"username": "Tester Updated",
		"gender":   "M",
	}

	tester := restit.Tester{
		BaseUrl: "http://localhost:8080/api.v1/users",
	}

	// -- Test Create --
	resp, err = tester.TestCreate(&userToCreate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}
	userId := result.Result[0].UserId // id of the created user

	// retrieve the user just created
	_, err = tester.TestRetrieveOne(fmt.Sprintf("%d", userId), &userToCreate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// -- Test Update --
	resp, err = tester.TestUpdate(fmt.Sprintf("%d", userId), &userToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// retrieve the user just updated
	_, err = tester.TestRetrieveOne(fmt.Sprintf("%d", userId), &userToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// -- Extended Test --
	// test: userSkill test
	err = testUserSkills(userId)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// -- Test Delete --
	// test: delete the user just created
	_, err = tester.TestDelete(fmt.Sprintf("%d", userId), &userToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	return
}
