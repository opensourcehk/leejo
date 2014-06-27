package main

import (
	"fmt"
	"github.com/yookoala/restit"
)

type SkillResp struct {
	Status string      `json:"status"`
	Result []UserSkill `json:"result"`
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
	ptr_a := a.(*UserSkill)
	ptr_b := b.(*UserSkill)
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

func testUserSkills(userId int64) (err error) {

	var resp SkillResp

	toCreate := UserSkill{
		UserId:    userId,
		SkillName: "Dummy Skill", // TODO: use uuid
	}
	toUpdate := UserSkill{
		UserId:    userId,
		SkillName: "Dummy Skill Updated",
	}

	test := restit.Rest("User",
		fmt.Sprintf("http://localhost:8080/api.v1/userSkills/%d", userId))

	// -- Test Create --
	test.Create(&toCreate).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toCreate).
		RunOrPanic()
	userSkillId := resp.Result[0].UserSkillId // id of the created user-skill

	test.Retrieve(fmt.Sprintf("%d", userSkillId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toCreate).
		RunOrPanic()

	// -- Test Update --
	test.Update(fmt.Sprintf("%d", userSkillId), &toUpdate).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		RunOrPanic()

	test.Retrieve(fmt.Sprintf("%d", userSkillId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		RunOrPanic()

	// -- Test Delete --
	test.Delete(fmt.Sprintf("%d", userSkillId)).
		WithResponseAs(&resp).
		ExpectResultCount(1).
		ExpectResultNth(0, &toUpdate).
		RunOrPanic()

	test.Retrieve(fmt.Sprintf("%d", userSkillId)).
		WithResponseAs(&resp).
		ExpectResultCount(0).
		RunOrPanic()

	return
}
