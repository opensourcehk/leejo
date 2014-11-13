package oauth2

import (
	"github.com/gourd/service"
	"github.com/gourd/session"
)

type ServiceProvider interface {
	ClientService(s session.Session) service.Service
	AuthService(s session.Session) service.Service
	AccessService(s session.Session) service.Service
	RefreshService(s session.Session) service.Service
}
