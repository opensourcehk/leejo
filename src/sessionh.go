package main

import (
	"net/http"
)

// definition of a session handler
type SessionHandler interface {
	// read the http request and returns a session
	Session(r *http.Request) (Session, error)
}

// implementat of a session handler
// that would fetch osin for authentication backend
type OsinSessionHandler struct {
}

func (h *OsinSessionHandler) Session(r *http.Request) (s Session, err error) {
	return
}
