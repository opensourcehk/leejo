package main

import (
	"data"
	"fmt"
	"github.com/yookoala/restit"
)

type InterestResp struct {
	Status string              `json:"status"`
	Result []data.UserInterest `json:"result"`
}

func (r *InterestResp) Count() int {
	return len(r.Result)
}

func (r *InterestResp) NthExists(n int) (err error) {
	if n < 0 || n > r.Count() {
		err = fmt.Errorf("Nth item (%d) not exist. Length = %d",
			n, len(r.Result))
	}
	return
}

func (r *InterestResp) NthValid(n int) (err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// test: the id should not be 0
	userInterestId := r.Result[n].UserInterestId
	if userInterestId == 0 {
		return fmt.Errorf("The user-interest has a UserInterestId = 0")
	}

	return
}

func (r *InterestResp) GetNth(n int) (o interface{}, err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// return nth pointer
	o = &r.Result[n]

	return

}

func (r *InterestResp) Match(a interface{}, b interface{}) (err error) {

	// check if the item match the payload
	ptr_a := a.(*data.UserInterest)
	ptr_b := b.(*data.UserInterest)
	if ptr_a.UserId != ptr_b.UserId {
		err = fmt.Errorf("UserId not match (\"%d\", \"%d\")",
			ptr_a.UserId, ptr_b.UserId)
		return
	} else if ptr_a.InterestName != ptr_b.InterestName {
		err = fmt.Errorf("InterestName not match (\"%s\", \"%s\")",
			ptr_a.InterestName, ptr_b.InterestName)
		return
	}

	return
}

func testUserInterests(token string, userId int64) (err error) {

	var resp InterestResp

	toCreate := data.UserInterest{
		UserId:       userId,
		InterestName: "Dummy Interest", // TODO: use uuid
	}
	toUpdate := data.UserInterest{
		UserId:       userId,
		InterestName: "Dummy Interest Updated",
	}

	test := restit.Rest("User",
		fmt.Sprintf("http://localhost:8080/api.v1/userInterests/%d", userId))

	// -- Test Create --
	test.Create(&toCreate).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toCreate).
		RunOrPanic()
	userInterestId := resp.Result[0].UserInterestId // id of the created user-interest

	test.Retrieve(fmt.Sprintf("%d", userInterestId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toCreate).
		RunOrPanic()

	// -- Test Update --
	test.Update(fmt.Sprintf("%d", userInterestId), &toUpdate).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		RunOrPanic()

	test.Retrieve(fmt.Sprintf("%d", userInterestId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		RunOrPanic()

	// -- Test Delete --
	test.Delete(fmt.Sprintf("%d", userInterestId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		RunOrPanic()

	test.Retrieve(fmt.Sprintf("%d", userInterestId)).
		WithResponseAs(&resp).
		ExpectResultCount(0).
		RunOrPanic()

	return
}
