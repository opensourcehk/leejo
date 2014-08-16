package session

import (
	"testing"
)

func TestScopes(t *testing.T) {
	t.Parallel()
	s := &Scopes{
		"whatever": true,
	}
	if !s.Has("whatever") {
		t.Errorf("Failed to obtain scope: %#v", *s)
	}
}

func TestScopes_Add(t *testing.T) {
	t.Parallel()
	s := &Scopes{
		"whatever": true,
	}
	if s.Has("foo") {
		t.Errorf("Incorrect assumtpion. Scope claimed to have something not in it", *s)
	}
	s.Add("foo")
	if !s.Has("foo") {
		t.Errorf("Failed to obtain added scope: %#v", *s)
	}
}

func TestScopes_Decode(t *testing.T) {
	t.Parallel()
	s := &Scopes{
		"whatever": true,
	}
	if s.Has("foo") {
		t.Errorf("Incorrect assumtpion. Scope claimed to have something not in it", *s)
	}
	s.Decode("foo,bar")
	if !s.Has("foo") || !s.Has("bar") {
		t.Errorf("Failed to obtain decoded scope: %#v", *s)
	}
}

func TestScopes_Encode(t *testing.T) {
	t.Parallel()
	s := &Scopes{
		"whatever": true,
		"foo":      true,
		"bar":      true,
	}
	result := s.Encode()
	if result != "bar,foo,whatever" &&
		result != "foo,whatever,bar" &&
		result != "whatever,bar,foo" &&
		result != "foo,bar,whatever" &&
		result != "bar,whatever,foo" &&
		result != "whatever,foo,bar" {
		t.Errorf("Encode gives unexpected result: %#v, %s", *s, result)
	}
}
