package rest

import (
	"github.com/gourd/session"
	"github.com/gourd/service"
)

// interface of helper that provides help to create
// a REST CURD interface
type Handler interface {

	// returns a pat readable regular expression
	// to listing endpoint
	BasePath() string

	// returns a pat readable regular expression
	// to individual entity
	EntityPath() string

	// allocate storage service for CURD operations
	Service(s session.Session) service.Service

	// translate an http request into a query context
	// i.e. key, parent key, query conditions, limit, offset and etc.
	Context(s session.Session) service.Context

	// check if the session allow
	// the kind of access to this object
	CheckAccess(string, session.Session, interface{}) error
}
