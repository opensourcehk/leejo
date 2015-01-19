package rest

import (
	"github.com/gourd/session"
	"log"
	"net/http"
)

// general rest error interface
type Error interface {
	GetCode() int
}

// default implementation of Error interface
type DefaultError struct {
	Status  string      `json:"status"`
	Code    int         `json:"code,omitempty"`
}

// implement Error interface GetCode method
func (e *DefaultError) GetCode() int {
	return e.Code
}

// write error to response
func WriteError(w http.ResponseWriter, sess session.Session, p Protocol, err error) {

	// log error first
	log.Printf("Internal Server Error: %s", err.Error())

	// get error by protocol
	pErr := p.WrapError(sess, err)

	// write header code
	if rErr, ok := pErr.(Error); ok {
		w.WriteHeader(rErr.GetCode())
	} else {
		w.WriteHeader(500)
	}

	// use protocol to encode
	p.NewEncoder(sess, w).Encode(pErr)

}
