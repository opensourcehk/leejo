package main

import (
	"github.com/yookoala/restit"
	"testing"
)

func TestSkillResp(t *testing.T) {
	t.Parallel()
	var i restit.Response
	i = &SkillResp{}
	i.Count()
}
