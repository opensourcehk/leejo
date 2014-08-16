package main

import (
	"github.com/gourd/service"
	"github.com/gourd/service/upperio"
	"leejo/data"
	"leejo/session"
	"upper.io/db"
)

// return a typical CURD service for User
type UserRest struct {
	Db       db.Database
	basePath string
	subPath  string
}

// check the session and see if it has the access
// that is required
func (h *UserRest) CheckAccess(access string, sess session.Session, ref interface{}) (err error) {
	u := sess.GetUser()
	if u == nil {
		err = service.Errorf(403, "Authentication Required")
	}
	return
}

// path which only include contextual information
// e.g. /api/user/123/emails
func (h *UserRest) BasePath() string {
	return h.basePath
}

// path which include entity specific information
// e.g. /api/user/123/emails/2
func (h *UserRest) EntityPath() string {
	return h.basePath + "/" + h.subPath
}

// allocate storage service for CURD operations of user
func (h *UserRest) Service(s session.Session) service.Service {
	// the content of service would be database specific
	// but the interface of service would be generic
	return &upperio.Service{
		Db:       h.Db,
		CollName: "leejo_user",
		IdSetterFunc: func(id upperio.Id, e service.EntityPtr) (err error) {
			u := e.(*data.User)
			u.UserId = id.(int64)
			return
		},
		CountFunc: func(el service.EntityListPtr) uint64 {
			l := el.(*[]data.User)
			return uint64(len(*l))
		},
		KeyCondFunc: func(k service.Key, pk service.ParentKey) service.Conds {
			cond := service.NewConds().Add("user_id", k)
			return cond
		},
		ListCondFunc: func(pk service.ParentKey) service.Conds {
			return service.NewConds().SetLimit(20)
		},
	}
}

// translate an http request into a query context
func (h *UserRest) Context(s session.Session) service.Context {
	c := service.NewConds().
		SetLimit(20).
		SetOffset(0)
	return &service.BasicContext{
		Key:       s.R().URL.Query().Get(":id"),
		ParentKey: nil,
		Conds:     c,
	}
}

// allocate an entity and return the address
func (h *UserRest) Entity() service.EntityPtr {
	return &data.User{}
}

// allocate a slice of entity and return the address
func (h *UserRest) EntityList() service.EntityListPtr {
	return &[]data.User{}
}

// get the length of a given pointer
func (h *UserRest) EntityListLen(p service.EntityListPtr) int {
	l := p.(*[]data.User)
	return len(*l)
}
