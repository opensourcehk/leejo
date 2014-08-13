package main

import (
	"github.com/yookoala/restit"
	"testing"
)

func TestInterestResp(t *testing.T) {
	t.Parallel()
	var i restit.Response
	i = &InterestResp{}
	i.Count()
}
