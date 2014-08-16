package session

import (
	"net/http"
)

type SessionUser interface{}

type Session interface {

	// obtain user
	GetUser() SessionUser

	// test if
	HasScope(string) bool

	// returns raw *http.Request
	R() *http.Request
}

// basic implementation of Session interface
type BasicSession struct {
	Request *http.Request
	Scopes  *Scopes
	User    SessionUser
}

func (a *BasicSession) GetUser() SessionUser {
	return a.User
}

func (a *BasicSession) HasScope(scope string) bool {
	if a.Scopes == nil {
		return false
	}
	return a.Scopes.Has(scope)
}

func (a *BasicSession) R() *http.Request {
	return a.Request
}
