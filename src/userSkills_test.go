package main

import (
	"testing"
)

func TestUserSkillRest(t *testing.T) {
	t.Parallel()
	var h PatRestHelper
	h = &UserSkillRest{}
	t.Logf("UserSkillRest as PatRestHelper: %#v", h)
}
