package main

import (
	"github.com/gourd/service"
	"github.com/gourd/service/upperio"
	"leejo/data"
	"leejo/session"
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
func (h *UserSkillRest) CheckAccess(access string, sess session.Session, ref interface{}) (err error) {
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
func (h *UserSkillRest) Service(s session.Session) service.Service {
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
		KeyCondFunc: func(k service.Key, pk service.ParentKey) service.Conds {
			c := service.NewConds().
				Add("user_id", pk).
				Add("user_skill_id", k)
			return c
		},
		ListCondFunc: func(pk service.ParentKey) service.Conds {
			c := service.NewConds().
				Add("user_id", pk).
				SetLimit(20)
			return c
		},
	}
}

// translate an http request into a query context
func (h *UserSkillRest) Context(s session.Session) service.Context {
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
func (h *UserSkillRest) Entity() service.EntityPtr {
	return &data.UserSkill{}
}

// allocate a slice of entity and return the address
func (h *UserSkillRest) EntityList() service.EntityListPtr {
	return &[]data.UserSkill{}
}
