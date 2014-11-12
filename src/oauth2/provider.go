package oauth2

import (
	"github.com/gourd/service"
	"github.com/gourd/session"
)

type ServiceProvider interface {
	Client(s session.Session) service.Service
	Auth(s session.Session) service.Service
	Access(s session.Session) service.Service
	Refresh(s session.Session) service.Service
}
