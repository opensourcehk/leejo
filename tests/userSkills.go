package main

import (
	"fmt"
	"github.com/jmcvetta/napping"
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

func (r *SkillResp) NthMatches(n int, comp *interface{}) (err error) {

	// check if the item exists
	err = r.NthExists(n)
	if err != nil {
		return
	}

	// check if the item match the payload
	userSkill := r.Result[n]
	cptr := (*comp).(*map[string]interface{})
	c := *cptr
	if userSkill.SkillName != c["skill_name"].(string) {
		err = fmt.Errorf("SkillName is \"%s\" (expected \"%s\")",
			userSkill.SkillName, c["skill_name"].(string))
		return
	} else if userSkill.UserId != c["user_id"].(int64) {
		err = fmt.Errorf("UserId is \"%d\" (expected \"%d\")",
			userSkill.UserId, c["user_id"].(int64))
		return
	}

	return
}

func testUserSkills(userId int64) (err error) {

	var result SkillResp
	var resp *napping.Response
	var userIdStr string

	userIdStr = fmt.Sprintf("%d", userId)

	skillToCreate := map[string]interface{}{
		"user_id":    userId,
		"skill_name": "Dummy Skill", // TODO: use uuid
	}
	skillToUpdate := map[string]interface{}{
		"user_id":    userId,
		"skill_name": "Dummy Skill Updated",
	}

	tester := restit.Tester{
		BaseUrl: "http://localhost:8080/api.v1/userSkills/" + userIdStr,
	}

	// -- Test Create --
	resp, err = tester.TestCreate(&skillToCreate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}
	userSkillId := result.Result[0].UserSkillId // id of the created user-skill

	// retrieve the user-skill just created
	_, err = tester.TestRetrieveOne(fmt.Sprintf("%d", userSkillId), &skillToCreate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// -- Test Update --
	resp, err = tester.TestUpdate(fmt.Sprintf("%d", userSkillId), &skillToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// retrieve the user-skill just updated
	_, err = tester.TestRetrieveOne(fmt.Sprintf("%d", userSkillId), &skillToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	// -- Test Delete --
	// test: delete the user-skill just created
	_, err = tester.TestDelete(fmt.Sprintf("%d", userSkillId), &skillToUpdate, &result)
	if err != nil {
		fmt.Printf("Raw: %s\n", resp.RawText())
		panic(err)
	}

	return
}
