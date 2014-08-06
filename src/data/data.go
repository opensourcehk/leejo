package data

type User struct {
	UserId   int64  `json:"user_id" db:"user_id,omitempty" form:"-"`
	Username string `json:"username" db:"username" form:"username"`
	Gender   string `json:"gender" db:"gender" form:"gender"`
}

type UserSkill struct {
	UserSkillId int64  `json:"user_skill_id" db:"user_skill_id,omitempty" form:"-"`
	UserId      int64  `json:"user_id" db:"user_id"`
	SkillName   string `json:"skill_name" db:"skill_name"`
}

type UserInterest struct {
	UserInterestId int64  `json:"user_interest_id" db:"user_interest_id,omitempty" form:"-"`
	UserId         int64  `json:"user_id" db:"user_id"`
	InterestName   string `json:"interest_name" db:"interest_name"`
}
