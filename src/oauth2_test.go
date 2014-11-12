package main

import (
	"leejo/oauth2"
	"testing"
)

func Test_oauth2Provider(t *testing.T) {
	var p oauth2.ServiceProvider
	p = &oauth2Provider{}
	t.Logf("oauth2Provider as oauth2.ServiceProvider: %#v", p)
}
