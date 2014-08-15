package session

import (
	"github.com/RangelReale/osin"
	"log"
	"net/http"
	"strings"
)

// definition of a session handler
type Handler interface {
	// read the http request and returns a session
	Session(r *http.Request) (Session, error)
}

// implementation of a session handler
// that would fetch osin for authentication backend
type OsinHandler struct {
	Storage osin.Storage
}

func (h *OsinHandler) Session(r *http.Request) (s Session, err error) {

	// this is the basic return anyway
	sps := make(Scopes)
	bs := BasicSession{
		Request: r,
		Scopes:  &sps,
	}
	s = &bs

	// read oauth2 token from header
	auth := r.Header.Get("Authorization")
	if auth == "" {
		log.Printf("Authorization empty")
		return
	}

	// parse the auth header
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
		log.Printf("incorrect Authorization header: %s", auth)
		return
	}

	// try to load access data with token
	a, err := h.Storage.LoadAccess(parts[1])
	if err != nil {
		return
	}

	// use the scope and user data loaded
	log.Printf("AccessData loaded: %#v", a)
	bs.User = a.UserData
	bs.Scopes.Decode(a.Scope)
	return
}
