package session

import (
	"strings"
)

// scope handler
type Scopes map[string]bool

// decode comma separated scope string
func (s *Scopes) Decode(scopesStr string) {
	strs := strings.Split(scopesStr, ",")
	for _, str := range strs {
		s.Add(strings.Trim(str, " \t\r\n"))
	}
}

// encode current content into comma separated string
func (s *Scopes) Encode() (scopesStr string) {
	arr := make([]string, len(*s))
	i := 0
	for str, _ := range *s {
		arr[i] = str
		i++
	}
	scopesStr = strings.Join(arr, ",")
	return
}

// add a scope
func (s *Scopes) Add(scope string) {
	(*s)[scope] = true
}

// delete a scope
func (s *Scopes) Del(scope string) {
	delete(*s, scope)
}

// check if a scope exists in current scopes
func (s *Scopes) Has(scope string) (ok bool) {
	_, ok = (*s)[scope]
	return
}
