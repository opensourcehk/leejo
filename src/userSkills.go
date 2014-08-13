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
type UserSkillRest struct {
	Db       db.Database
	basePath string
	subPath  string
}

// check the session and see if it has the access
// that is required
func (h *UserSkillRest) CheckAccess(access string, sh SessionHandler, ref interface{}) (err error) {
	return
}

// path which only include contextual information
// e.g. /api/user/123/emails
func (h *UserSkillRest) BasePath() string {
	return h.basePath
}

// path which include entity specific information
// e.g. /api/user/123/emails/2
func (h *UserSkillRest) EntityPath() string {
	return h.basePath + "/" + h.subPath
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
		KeyCondFunc: func(k service.Key, pk service.ParentKey) db.Cond {
			return db.Cond{
				"user_id":       pk,
				"user_skill_id": k,
			}
		},
		ParentCondFunc: func(pk service.ParentKey) db.Cond {
			return db.Cond{
				"user_skill_id": pk,
			}
		},
	}
}

// translate an http request into a query context
func (h *UserSkillRest) Context(r *http.Request) service.Context {
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
func (h *UserSkillRest) Entity() service.EntityPtr {
	return &data.UserSkill{}
}

// allocate a slice of entity and return the address
func (h *UserSkillRest) EntityList() service.EntityListPtr {
	return &[]data.UserSkill{}
}

// create user CURD interface with pat
func bindUserSkills(path string, sh SessionHandler, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr

	h := UserSkillRest{
		Db:       sess,
		basePath: path,
		subPath:  "{id:[0-9]+}",
	}
	RestOnPat(&h, sh, r)
}
