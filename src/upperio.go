package main

import (
	"leejo/service"
	"upper.io/db"
)

// receives the serial / id from upper.io
type EntityId interface{}

// Generic upper.io CURD service
// Implements Service interface in leejo/service
type UpperIoService struct {
	Db             db.Database
	CollName       string
	ApplyIdFunc    func(EntityId, service.EntityPtr) (err error)
	KeyCondFunc    func(service.KeyPtr) db.Cond
	ParentCondFunc func(service.ParentKeyPtr) db.Cond
}

func (s *UpperIoService) ApplyId(id EntityId, e service.EntityPtr) (err error) {
	return s.ApplyIdFunc(id, e)
}

func (s *UpperIoService) KeyCond(k service.KeyPtr) db.Cond {
	return s.KeyCondFunc(k)
}

func (s *UpperIoService) ParentCond(pk service.ParentKeyPtr) db.Cond {
	return s.ParentCondFunc(pk)
}

func (s *UpperIoService) Create(c service.Context, e service.EntityPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	// add the user to user collection
	id, err := coll.Append(e)
	if err != nil {
		return
	}

	// apply the serial key to the
	err = s.ApplyId(id, e)
	return
}

func (s *UpperIoService) List(c service.Context, el service.EntityListPtr) (err error) {

	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	// retrieve all users
	res := coll.Find(s.ParentCond(c.ParentKey))
	// TODO: also work with c.Cond for ListCond (limit and offset)
	err = res.All(el)
	return
}

func (s *UpperIoService) Retrieve(c service.Context, el service.EntityListPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	// retrieve all users of id(s)
	res := coll.Find(s.KeyCond(c.Key))
	err = res.All(el)
	return
}

func (s *UpperIoService) Update(c service.Context, e service.EntityPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	res := coll.Find(s.KeyCond(c.Key))

	// update the user
	err = res.Update(e)
	return
}

func (s *UpperIoService) Delete(c service.Context) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	res := coll.Find(s.KeyCond(c.Key))
	if err != nil {
		return
	}

	err = res.Remove()
	return
}
