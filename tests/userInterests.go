package main

import (
	"fmt"
	"github.com/jmcvetta/napping"
	"github.com/yookoala/restit"
)

type InterestResp struct {
	Status string         `json:"status"`
	Result []UserInterest `json:"result"`
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

func (r *InterestResp) NthMatches(n int, comp *interface{}) (err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// check if the item match the payload
	userInterest := r.Result[n]
	cptr := (*comp).(*map[string]interface{})
	c := *cptr
	if userInterest.InterestName != c["interest_name"].(string) {
		err = fmt.Errorf("InterestName is \"%s\" (expected \"%s\")",
			userInterest.InterestName, c["interest_name"].(string))
		return
	} else if userInterest.UserId != c["user_id"].(int64) {
		err = fmt.Errorf("UserId is \"%d\" (expected \"%d\")",
			userInterest.UserId, c["user_id"].(int64))
		return
	}

	return
}

func testUserInterests(userId int64) (err error) {

	var result InterestResp
	var resp *napping.Response
	var userIdStr string

	userIdStr = fmt.Sprintf("%d", userId)

	interestToCreate := map[string]interface{}{
		"user_id":       userId,
		"interest_name": "Dummy Interest", // TODO: use uuid
	}
	interestToUpdate := map[string]interface{}{
		"user_id":       userId,
		"interest_name": "Dummy Interest Updated",
	}

	tester := restit.Tester{
		BaseUrl: "http://localhost:8080/api.v1/userInterests/" + userIdStr,
	}

	// -- Test Create --
	resp, err = tester.TestCreate(&interestToCreate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}
	userInterestId := result.Result[0].UserInterestId // id of the created user-interest

	// retrieve the user-interest just created
	_, err = tester.TestRetrieveOne(fmt.Sprintf("%d", userInterestId), &interestToCreate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// -- Test Update --
	resp, err = tester.TestUpdate(fmt.Sprintf("%d", userInterestId), &interestToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// retrieve the user-interest just updated
	_, err = tester.TestRetrieveOne(fmt.Sprintf("%d", userInterestId), &interestToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// -- Test Delete --
	// test: delete the user-interest just created
	_, err = tester.TestDelete(fmt.Sprintf("%d", userInterestId), &interestToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	return
}
