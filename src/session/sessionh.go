package session

import (
	"github.com/RangelReale/osin"
	"net/http"
)

// definition of a session handler
type SessionHandler interface {
	// read the http request and returns a session
	Session(r *http.Request) (Session, error)
}

// implementation of a session handler
// that would fetch osin for authentication backend
type OsinSessionHandler struct {
	Storage osin.Storage
}

func (h *OsinSessionHandler) Session(r *http.Request) (s Session, err error) {
	auth := r.Header.Get("Authorization")
	_ = auth

	s = &BasicSession{
		Request: r,
	}
	return
}
