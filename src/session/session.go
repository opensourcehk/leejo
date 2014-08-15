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
	return a.Scopes.Has(scope)
}

func (a *BasicSession) R() *http.Request {
	return a.Request
}

// scope handler
type Scopes map[string]bool

// decode comma separated scope string
func (s *Scopes) Decode(scopesStr string) {
}

// encode current content into comma separated string
func (s *Scopes) Encode() (scopesStr string) {
	return
}

// add a scope
func (s *Scopes) Add(scope string) {
}

// delete a scope
func (s *Scopes) Del(scope string) {
}

// check if a scope exists in current scopes
func (s *Scopes) Has(scope string) bool {
	return false
}
