package session

import (
	"testing"
)

func TestBasicSession(t *testing.T) {
	t.Parallel()
	var s Session
	s = &BasicSession{}
	s.HasScope("whatever")
}
