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
type UserSkillRest struct {
	Db db.Database
}

// path which include entity specific information
func (h *UserSkillRest) SubPath() string {
	return "{id:[0-9]+}"
}

// allocate storage service for CURD operations of user
func (h *UserSkillRest) Service(r *http.Request) service.Service {
	// the content of service would be database specific
	// but the interface of service would be generic
	return &upperio.Service{
		Db:       h.Db,
		CollName: "leejo_user_skill",
		IdSetterFunc: func(id upperio.Id, e service.EntityPtr) (err error) {
			u := e.(*data.UserSkill)
			u.UserSkillId = id.(int64)
			return
		},
		CountFunc: func(el service.EntityListPtr) uint64 {
			l := el.(*[]data.UserSkill)
			return uint64(len(*l))
		},
		KeyCondFunc: func(c service.Context) db.Cond {
			return db.Cond{
				"user_id":       c.Get(":user_id"),
				"user_skill_id": c.Get(":id"),
			}
		},
		ParentCondFunc: func(c service.Context) db.Cond {
			return db.Cond{
				"user_id": c.Get(":user_id"),
			}
		},
	}
}

// translate an http request into a query context
func (h *UserSkillRest) Context(r *http.Request) service.Context {
	return &service.BasicContext{
		Values: r.URL.Query(),
		Cond: &service.BasicListCond{
			Limit:  20,
			Offset: 0,
		},
	}
}

// allocate an entity and return the address
func (h *UserSkillRest) Entity() service.EntityPtr {
	return &data.UserSkill{}
}

// allocate a slice of entity and return the address
func (h *UserSkillRest) EntityList() service.EntityListPtr {
	return &[]data.UserSkill{}
}

// create user CURD interface with pat
func bindUserSkills(path string, osinServer *osin.Server, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr

	h := UserSkillRest{
		Db: sess,
	}
	RestOnPat(path, &h, r)
}
