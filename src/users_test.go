package main

import (
	"testing"
)

func TestUserRest(t *testing.T) {
	t.Parallel()
	var h PatRestHelper
	h = &UserRest{}
	t.Logf("UserRest as PatRestHelper: %#v", h)
}
