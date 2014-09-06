package rest

import (
	"github.com/gourd/session"
	"net/http"
)

// general rest error interface
type Response interface {
	GetCode() int
}

// default implementation of response implements Error
type DefaultResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

// implements GetCode method of Response interface
func (r *DefaultResponse) GetCode() int {
	return r.Code
}

// write response to response
func WriteResponse(w http.ResponseWriter, sess session.Session, p Protocol, r interface{}, c int) {

	// get response by protocol
	pr := p.Response(sess, r)

	// try to retrieve status code from protocol
	if prr, ok := pr.(Response); ok {
		if pc := prr.GetCode(); pc != 0 {
			c = pc
		}
	}

	// write header with status code
	w.WriteHeader(c)

	// use protocol to encode
	p.NewEncoder(sess, w).Encode(pr)

}
