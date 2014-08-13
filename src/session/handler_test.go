package session

import (
	"testing"
)

func TestOsinSessionHandler(t *testing.T) {
	t.Parallel()
	var h SessionHandler
	h = &OsinSessionHandler{}
	t.Logf("Obtained session hander: %#v", h)
}
