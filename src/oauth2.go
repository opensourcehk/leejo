package main

import (
	"github.com/gourd/service"
	"github.com/gourd/service/upperio"
	"github.com/gourd/session"
	"leejo/data"
	"upper.io/db"
)

type oauth2Provider struct {
	Db db.Database
}

// allocate storage service for CURD operations of user
func (p *oauth2Provider) ClientService(s session.Session) service.Service {
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
				Add("id", k)
			return c
		},
		ListCondFunc: func(pk service.ParentKey) service.Conds {
			c := service.NewConds().
				Add("id", pk).
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

// allocate storage service for CURD operations of user
func (p *oauth2Provider) AuthService(s session.Session) service.Service {
	// the content of service would be database specific
	// but the interface of service would be generic
	return &upperio.Service{
		Db:       p.Db,
		CollName: "leejo_api_authdata",
		IdSetterFunc: func(id upperio.Id, e service.EntityPtr) (err error) {
			d := e.(*data.ApiAuthData)
			d.Id = id.(int64)
			return
		},
		CountFunc: func(el service.EntityListPtr) uint64 {
			l := el.(*[]data.ApiAuthData)
			return uint64(len(*l))
		},
		KeyCondFunc: func(k service.Key, pk service.ParentKey) service.Conds {
			c := service.NewConds().
				Add("code", k)
			return c
		},
		ListCondFunc: func(pk service.ParentKey) service.Conds {
			c := service.NewConds().
				Add("code", pk)
			return c
		},
		EntityFunc: func() service.EntityPtr {
			return &data.ApiAuthData{}
		},
		EntityListFunc: func() service.EntityListPtr {
			return &[]data.ApiAuthData{}
		},
		LenFunc: func(p service.EntityListPtr) int64 {
			l := p.(*[]data.ApiAuthData)
			return int64(len(*l))
		},
	}
}

// allocate storage service for CURD operations of user
func (p *oauth2Provider) AccessService(s session.Session) service.Service {
	// the content of service would be database specific
	// but the interface of service would be generic
	return &upperio.Service{
	// placeholder only, for now
	}
}

// allocate storage service for CURD operations of user
func (p *oauth2Provider) RefreshService(s session.Session) service.Service {
	// the content of service would be database specific
	// but the interface of service would be generic
	return &upperio.Service{
	// placeholder only, for now
	}
}
