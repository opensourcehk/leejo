package session

import (
	"testing"
)

func TestOsinHandler(t *testing.T) {
	t.Parallel()
	var h Handler
	h = &OsinHandler{}
	t.Logf("Obtained session hander: %#v", h)
}
