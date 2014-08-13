package main

import (
	"github.com/RangelReale/osin"
	"github.com/gorilla/pat"
	"github.com/gourd/service"
	"github.com/gourd/service/upperio"
	"leejo/data"
	"net/http"
	"upper.io/db"
)

// return a typical CURD service for User
type UserInterestRest struct {
	Db       db.Database
	basePath string
	subPath  string
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
func (h *UserInterestRest) Service(r *http.Request) service.Service {
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
		KeyCondFunc: func(c service.Context) db.Cond {
			return db.Cond{
				"user_id":          c.GetParentKey(),
				"user_interest_id": c.GetKey(),
			}
		},
		ParentCondFunc: func(c service.Context) db.Cond {
			return db.Cond{
				"user_id": c.GetParentKey(),
			}
		},
	}
}

// translate an http request into a query context
func (h *UserInterestRest) Context(r *http.Request) service.Context {
	q := r.URL.Query()
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

// create user CURD interface with pat
func bindUserInterests(path string, osinServer *osin.Server, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr

	h := UserInterestRest{
		Db:       sess,
		basePath: path,
		subPath:  "{id:[0-9]+}",
	}
	RestOnPat(&h, r)
}
