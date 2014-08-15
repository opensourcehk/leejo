package data

import (
	"leejo/oauth2"
	"testing"
)

func TestUser(t *testing.T) {
	t.Parallel()

	var u oauth2.User
	u = &User{}
	t.Logf("Can cast data.User to oauth2.User: %#v", u)
}

func TestUser_Casting(t *testing.T) {
	t.Parallel()

	var u oauth2.User
	var i interface{}
	var ok bool
	i = &User{}
	if u, ok = i.(oauth2.User); ok {
		t.Logf("Can cast data.User, through empty interface, to oauth2.User: %#v", u)
	} else {
		t.Errorf("Cannot cast data.User, through empty interface, to oauth2.User: %#v", u)
	}
}
