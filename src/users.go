package main

import (
	"github.com/RangelReale/osin"
	"github.com/gorilla/pat"
	"leejo/data"
	"leejo/service"
	"net/http"
	"upper.io/db"
)

// return a typical CURD service for User
type UserRest struct {
	Db db.Database
}

// path which include entity specific information
func (h *UserRest) SubPath() string {
	return "{id:[0-9]+}"
}

func (h *UserRest) Service(r *http.Request) service.Service {
	// the content of service would be database specific
	// but the interface of service would be generic
	return &UpperIoService{
		Db:       h.Db,
		CollName: "leejo_user",
		ApplyIdFunc: func(id EntityId, e service.EntityPtr) (err error) {
			u := e.(*data.User)
			u.UserId = id.(int64)
			return
		},
		KeyCondFunc: func(k service.KeyPtr) db.Cond {
			return db.Cond{"user_id": k}
		},
		ParentCondFunc: func(pk service.ParentKeyPtr) db.Cond {
			return db.Cond{}
		},
	}
}

// translate an http request into a query context
func (h *UserRest) Context(r *http.Request) service.Context {
	return service.Context{
		Key:       r.URL.Query().Get(":id"),
		ParentKey: nil,
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
func bindUser(path string, osinServer *osin.Server, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr

	h := UserRest{
		Db: sess,
	}
	RestOnPat(path, &h, r)
}
