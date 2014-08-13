package main

func ObtainSession(code string) (s Session, err error) {
	return
}

type SessionUser interface{}

type Session interface {

	// test if
	HasScope(string) bool

	// obtain user
	GetUser() SessionUser
}

// basic implementation of Session interface
type BasicSession struct {
	Scopes *Scopes
	User   SessionUser
}

func (a *BasicSession) HasScope(scope string) bool {
	return a.Scopes.Has(scope)
}

func (a *BasicSession) GetUser() (u SessionUser) {
	return
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
