package oauth2

import (
	"github.com/gourd/service"
	"github.com/gourd/session"
)

// interface of helper that provides help to
// provide service with accordance to a given session
type ServiceProvider interface {

	// allocate storage service for CURD operations
	Service(s session.Session) service.Service
}
