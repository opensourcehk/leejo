package data

import (
	"testing"
)

func TestUser(t *testing.T) {
	t.Parallel()

	var u ApiUser
	u = &User{}
	t.Logf("Can cast data.User to oauth2.User: %#v", u)
}

func TestUser_Casting(t *testing.T) {
	t.Parallel()

	var u ApiUser
	var i interface{}
	var ok bool
	i = &User{}
	if u, ok = i.(ApiUser); ok {
		t.Logf("Can cast data.User, through empty interface, to oauth2.User: %#v", u)
	} else {
		t.Errorf("Cannot cast data.User, through empty interface, to oauth2.User: %#v", u)
	}
}
