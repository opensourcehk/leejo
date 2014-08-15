package oauth2

import (
	"testing"
)

func Test_apiAuthUser(t *testing.T) {
	var u User
	u = &apiAuthUser{}
	t.Logf("apiAuthUser implements User: %#v", u)
}
