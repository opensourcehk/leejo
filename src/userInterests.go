package main

import (
	"github.com/gourd/service"
	"github.com/gourd/service/upperio"
	"leejo/data"
	"leejo/session"
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
		KeyCondFunc: func(k service.Key, pk service.ParentKey) db.Cond {
			return db.Cond{
				"user_id":          pk,
				"user_interest_id": k,
			}
		},
		ParentCondFunc: func(pk service.ParentKey) db.Cond {
			return db.Cond{
				"user_id": pk,
			}
		},
	}
}

// translate an http request into a query context
func (h *UserInterestRest) Context(s session.Session) service.Context {
	q := s.R().URL.Query()
	return &service.BasicContext{
		Key:       q.Get(":id"),
		ParentKey: q.Get(":user_id"),
		Values:    q,
		Cond: &service.BasicListCond{
			Limit:  20,
			Offset: 0,
		},
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
