package main

import (
	"fmt"
	"github.com/yookoala/restit"
	"leejo/data"
)

type UserResp struct {
	Status string      `json:"status"`
	Result []data.User `json:"result"`
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

func (r *UserResp) GetNth(n int) (uo interface{}, err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// return nth pointer
	uo = &r.Result[n]

	return

}

func (r *UserResp) Match(a interface{}, b interface{}) (err error) {

	// check if the item match the payload
	ptr_a := a.(*data.User)
	ptr_b := b.(*data.User)
	if ptr_a.Username != ptr_b.Username {
		err = fmt.Errorf("Username not match (\"%s\", \"%s\")",
			ptr_a.Username, ptr_b.Username)
		return
	} else if ptr_a.Gender != ptr_b.Gender {
		err = fmt.Errorf("Gender not match (\"%s\", \"%s\")",
			ptr_a.Gender, ptr_b.Gender)
		return
	}

	return
}

func testUser(token string) (err error) {

	var resp UserResp

	userToCreate := data.User{
		Username: "Tester", // TODO: use uuid
		Gender:   "F",
	}
	userToUpdate := data.User{
		Username: "Tester Updated",
		Gender:   "M",
	}

	test := restit.Rest("User", "http://localhost:8080/api.v1/users")

	// -- Test Create without proper token --
	test.Create(&userToCreate).
		ExpectStatus(403).
		RunOrPanic()

	// -- Test Create --
	test.Create(&userToCreate).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &userToCreate).
		RunOrPanic()
	userId := resp.Result[0].UserId // id of the created user

	test.Retrieve(fmt.Sprintf("%d", userId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &userToCreate).
		RunOrPanic()

	// -- Test Update --
	test.Update(fmt.Sprintf("%d", userId), &userToUpdate).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &userToUpdate).
		RunOrPanic()

	test.Retrieve(fmt.Sprintf("%d", userId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &userToUpdate).
		RunOrPanic()

	// -- Extended Test --
	// test: userSkill test
	err = testUserSkills(token, userId)
	if err != nil {
		return
	}

	// test: userInterest test
	err = testUserInterests(token, userId)
	if err != nil {
		return
	}

	// -- Test Delete --
	test.Delete(fmt.Sprintf("%d", userId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &userToUpdate).
		RunOrPanic()

	test.Retrieve(fmt.Sprintf("%d", userId)).
		WithResponseAs(&resp).
		ExpectResultCount(0).
		RunOrPanic()

	return
}
