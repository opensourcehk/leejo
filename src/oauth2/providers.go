package oauth2

import (
	"github.com/gourd/service"
	"github.com/gourd/service/upperio"
	"github.com/gourd/session"
	"leejo/data"
	"upper.io/db"
)

type ClientProvider struct {
	Db db.Database
}

// allocate storage service for CURD operations of user
func (p *ClientProvider) Service(s session.Session) service.Service {
	// the content of service would be database specific
	// but the interface of service would be generic
	return &upperio.Service{
		Db:       p.Db,
		CollName: "leejo_api_client",
		IdSetterFunc: func(id upperio.Id, e service.EntityPtr) (err error) {
			c := e.(*data.ApiClient)
			c.Id = id.(string)
			return
		},
		CountFunc: func(el service.EntityListPtr) uint64 {
			l := el.(*[]data.ApiClient)
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
			return &data.ApiClient{}
		},
		EntityListFunc: func() service.EntityListPtr {
			return &[]data.ApiClient{}
		},
		LenFunc: func(p service.EntityListPtr) int64 {
			l := p.(*[]data.ApiClient)
			return int64(len(*l))
		},
	}
}

