package main

import (
	"leejo/rest"
	"testing"
)

func TestUserSkillRest(t *testing.T) {
	t.Parallel()
	var h rest.Handler
	h = &UserSkillRest{}
	t.Logf("UserSkillRest as PatRestHelper: %#v", h)
}
