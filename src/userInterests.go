package main

import (
	"github.com/gourd/service"
	"github.com/gourd/service/upperio"
	"github.com/gourd/session"
	"leejo/data"
	"upper.io/db"
)

// return a typical CURD service for User
type UserInterestRest struct {
	Db       db.Database
	basePath string
	subPath  string
}

// check the session and see if it has the access
// that is required
func (h *UserInterestRest) CheckAccess(access string, sess session.Session, ref interface{}) (err error) {
	u := sess.GetUser()
	if u == nil {
		err = service.Errorf(403, "Authentication Required")
	}
	if !sess.HasScope("user_interests") {
		err = service.Errorf(403, "Lack the Required Scope")
	}
	return
}

// path which only include contextual information
// e.g. /api/user/123/emails
func (h *UserInterestRest) BasePath() string {
	return h.basePath
}

// path which include entity specific information
// e.g. /api/user/123/emails/2
func (h *UserInterestRest) EntityPath() string {
	return h.basePath + "/" + h.subPath
}

// allocate storage service for CURD operations of user
func (h *UserInterestRest) Service(s session.Session) service.Service {
	// the content of service would be database specific
	// but the interface of service would be generic
	return &upperio.Service{
		Db:       h.Db,
		CollName: "leejo_user_interest",
		IdSetterFunc: func(id upperio.Id, e service.EntityPtr) (err error) {
			u := e.(*data.UserInterest)
			u.UserInterestId = id.(int64)
			return
		},
		CountFunc: func(el service.EntityListPtr) uint64 {
			l := el.(*[]data.UserInterest)
			return uint64(len(*l))
		},
		KeyCondFunc: func(k service.Key, pk service.ParentKey) service.Conds {
			c := service.NewConds().
				Add("user_id", pk).
				Add("user_interest_id", k)
			return c
		},
		ListCondFunc: func(pk service.ParentKey) service.Conds {
			c := service.NewConds().
				Add("user_id", pk).
				SetLimit(20)
			return c
		},
		EntityFunc: func() service.EntityPtr {
			return &data.UserInterest{}
		},
		EntityListFunc: func() service.EntityListPtr {
			return &[]data.UserInterest{}
		},
		LenFunc: func(p service.EntityListPtr) int64 {
			l := p.(*[]data.UserInterest)
			return int64(len(*l))
		},
	}
}

// translate an http request into a query context
func (h *UserInterestRest) Context(s session.Session) service.Context {
	q := s.R().URL.Query()
	c := service.NewConds().
		SetLimit(20).
		SetOffset(0)
	return &service.BasicContext{
		Key:       q.Get(":id"),
		ParentKey: q.Get(":user_id"),
		Conds:     c,
	}
}

// allocate an entity and return the address
func (h *UserInterestRest) Entity() service.EntityPtr {
	return &data.UserInterest{}
}

// allocate a slice of entity and return the address
func (h *UserInterestRest) EntityList() service.EntityListPtr {
	return &[]data.UserInterest{}
}

// get the length of a given pointer
func (h *UserInterestRest) EntityListLen(p service.EntityListPtr) int {
	l := p.(*[]data.UserInterest)
	return len(*l)
}
