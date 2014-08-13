package main

import (
	"github.com/gorilla/pat"
	"github.com/gourd/service"
	"github.com/gourd/service/upperio"
	"leejo/data"
	"net/http"
	"upper.io/db"
)

// return a typical CURD service for User
type UserRest struct {
	Db       db.Database
	basePath string
	subPath  string
}

// return the scope requires for
// a certain kind of access type
func (h *UserRest) ScopeOf(access string) string {
	return "user"
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
func (h *UserRest) Service(r *http.Request) service.Service {
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
		KeyCondFunc: func(k service.Key, pk service.ParentKey) db.Cond {
			return db.Cond{"user_id": k}
		},
		ParentCondFunc: func(pk service.ParentKey) db.Cond {
			return db.Cond{}
		},
	}
}

// translate an http request into a query context
func (h *UserRest) Context(r *http.Request) service.Context {
	q := r.URL.Query()
	return &service.BasicContext{
		Key:       q.Get(":id"),
		ParentKey: nil,
		Values:    q,
		Cond: &service.BasicListCond{
			Limit:  20,
			Offset: 0,
		},
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

// create user CURD interface with pat
func bindUser(path string, sh SessionHandler, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr

	h := UserRest{
		Db:       sess,
		basePath: path,
		subPath:  "{id:[0-9]+}",
	}
	RestOnPat(&h, r)
}
