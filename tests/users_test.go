package main

import (
	"github.com/yookoala/restit"
	"testing"
)

func TestUserResp(t *testing.T) {
	t.Parallel()
	var i restit.Response
	i = &UserResp{}
	i.Count()
}
