package main

type User struct {
	UserId   uint32 `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	Gender   string `json:"gender" db:"gender"`
}

type Skill struct {
	UserId    uint32 `json:"user_id" db:"user_id"`
	SkillName string `json:"skill_name" db:"skill_name"`
}

type Interest struct {
	UserId       uint32 `json:user_id db:"user_id"`
	InterestName string `json:interest_name db:"interest_name"`
}
