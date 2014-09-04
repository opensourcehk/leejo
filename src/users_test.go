package main

import (
	"leejo/rest"
	"testing"
)

func TestUserRest(t *testing.T) {
	t.Parallel()
	var h rest.Handler
	h = &UserRest{}
	t.Logf("UserRest as PatRestHelper: %#v", h)
}
