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

func TestBasicSession_Scope(t *testing.T) {
	t.Parallel()
	s := BasicSession{
		Scopes: &Scopes{
			"whatever": true,
		},
	}
	if !s.HasScope("whatever") {
		t.Errorf("Failed to obtain scope: %#v", s.Scopes)
	}
}
