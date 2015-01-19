package data

import (
	"testing"
)

func Test_ApiAuthUser(t *testing.T) {
	var u ApiUser
	u = &ApiAuthUser{}
	t.Logf("ApiAuthUser implements User: %#v", u)
}
