package main

import (
	"github.com/opensourcehk/leejo/lib/rest"
	"testing"
)

func TestUserInterestRest(t *testing.T) {
	t.Parallel()
	var h rest.Handler
	h = &UserInterestRest{}
	t.Logf("UserInterestRest as PatRestHelper: %#v", h)
}
