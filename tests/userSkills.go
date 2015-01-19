package main

import (
	"fmt"
	"github.com/opensourcehk/leejo/lib/data"
	"github.com/yookoala/restit"
)

type SkillResp struct {
	Status string           `json:"status"`
	Result []data.UserSkill `json:"result"`
}

func (r *SkillResp) Count() int {
	return len(r.Result)
}

func (r *SkillResp) NthExists(n int) (err error) {
	if n < 0 || n > r.Count() {
		err = fmt.Errorf("Nth item (%d) not exist. Length = %d",
			n, len(r.Result))
	}
	return
}

func (r *SkillResp) NthValid(n int) (err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// test: the id should not be 0
	userSkillId := r.Result[n].UserSkillId
	if userSkillId == 0 {
		return fmt.Errorf("The user-skill has a UserSkillId = 0")
	}

	return
}

func (r *SkillResp) GetNth(n int) (o interface{}, err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// return nth pointer
	o = &r.Result[n]

	return

}

func (r *SkillResp) Match(a interface{}, b interface{}) (err error) {

	// check if the item match the payload
	ptr_a := a.(*data.UserSkill)
	ptr_b := b.(*data.UserSkill)
	if ptr_a.UserId != ptr_b.UserId {
		err = fmt.Errorf("UserId not match (\"%d\", \"%d\")",
			ptr_a.UserId, ptr_b.UserId)
		return
	} else if ptr_a.SkillName != ptr_b.SkillName {
		err = fmt.Errorf("SkillName not match (\"%s\", \"%s\")",
			ptr_a.SkillName, ptr_b.SkillName)
		return
	}

	return
}

func (r *SkillResp) Reset() {
	r.Result = make([]data.UserSkill, 0)
}

func testUserSkills(token string, userId int64) (err error) {

	var resp SkillResp

	toCreate := data.UserSkill{
		UserId:    userId,
		SkillName: "Dummy Skill", // TODO: use uuid
	}
	toUpdate := data.UserSkill{
		UserId:    userId,
		SkillName: "Dummy Skill Updated",
	}

	test := restit.Rest("Skill",
		fmt.Sprintf("http://localhost:8080/api.v1/userSkills/%d", userId))

	// -- Test Create --
	test.Create(&toCreate).
		AddHeader("Authorization", "Bearer "+token).
		WithResponseAs(&resp).
		ExpectStatus(201). // Created
		ExpectResultCount(1).
		ExpectResultNth(0, &toCreate).
		RunOrPanic()
	userSkillId := resp.Result[0].UserSkillId // id of the created user-skill

	test.Retrieve(fmt.Sprintf("%d", userSkillId)).
		AddHeader("Authorization", "Bearer "+token).
		WithResponseAs(&resp).
		ExpectStatus(200). // Success
		ExpectResultCount(1).
		ExpectResultNth(0, &toCreate).
		RunOrPanic()

	// -- Test Update --
	test.Update(fmt.Sprintf("%d", userSkillId), &toUpdate).
		AddHeader("Authorization", "Bearer "+token).
		WithResponseAs(&resp).
		ExpectStatus(200). // Success
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		RunOrPanic()

	test.Retrieve(fmt.Sprintf("%d", userSkillId)).
		AddHeader("Authorization", "Bearer "+token).
		WithResponseAs(&resp).
		ExpectStatus(200). // Success
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		RunOrPanic()

	// -- Test Delete --
	test.Delete(fmt.Sprintf("%d", userSkillId)).
		AddHeader("Authorization", "Bearer "+token).
		WithResponseAs(&resp).
		ExpectStatus(404). // Not Found
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		Run()

	test.Retrieve(fmt.Sprintf("%d", userSkillId)).
		AddHeader("Authorization", "Bearer "+token).
		ExpectStatus(404). // Not Found
		WithResponseAs(&resp).
		ExpectResultCount(0).
		RunOrPanic()

	return
}
