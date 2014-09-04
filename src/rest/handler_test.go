package rest

import (
	"github.com/gourd/session"
	"github.com/gourd/service"
	"testing"
)

// interface of helper that provides help to create
// a REST CURD interface
type testHandler struct {
}

func (h *testHandler) BasePath() string {
	return ""
}

func (h *testHandler) EntityPath() string {
	return ""
}

func (h *testHandler) Service(session.Session) service.Service {
	return nil
}

func (h *testHandler) Context(session.Session) service.Context {
	return nil
}

func (h *testHandler) CheckAccess(string, session.Session, interface{}) error {
	return nil
}


func TestTestHandler(t *testing.T) {
	var h Handler
	h = &testHandler{}
	t.Log("testHandler implements Handler: %#v", h)
}
