package main

import (
	"testing"
)

func TestUserInterestRest(t *testing.T) {
	t.Parallel()
	var h PatRestHelper
	h = &UserInterestRest{}
	t.Logf("UserInterestRest as PatRestHelper: %#v", h)
}
