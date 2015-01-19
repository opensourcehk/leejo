package main

import (
	"github.com/opensourcehk/leejo/lib/rest"
	"testing"
)

func TestUserSkillRest(t *testing.T) {
	t.Parallel()
	var h rest.Handler
	h = &UserSkillRest{}
	t.Logf("UserSkillRest as PatRestHelper: %#v", h)
}
