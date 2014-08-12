package upperio

import (
	"leejo/service"
	"upper.io/db"
)

// receives the serial / id from upper.io
// used in Append operations
type Id interface{}

// Generic upper.io CURD service
// Implements Service interface in leejo/service
type Service struct {
	Db             db.Database
	CollName       string
	IdSetterFunc   func(Id, service.EntityPtr) (err error)
	KeyCondFunc    func(service.Context) db.Cond
	ParentCondFunc func(service.Context) db.Cond
}

// upperio specific method
// apply id to the just created entity
// used in Create method
func (s *Service) SetId(id Id, e service.EntityPtr) (err error) {
	return s.IdSetterFunc(id, e)
}

// upperio specific method
// translate key into upper.io condition
// used in Retrieve method
func (s *Service) KeyCond(c service.Context) db.Cond {
	return s.KeyCondFunc(c)
}

// upperio specific method
// translate parent key into upper.io condition
// used in List method
func (s *Service) ParentCond(c service.Context) db.Cond {
	return s.ParentCondFunc(c)
}

// implements Create method of Service
func (s *Service) Create(c service.Context, e service.EntityPtr) (err error) {
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
	err = s.SetId(id, e)
	return
}

// implements List method of Service
func (s *Service) List(c service.Context, el service.EntityListPtr) (err error) {

	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	// retrieve all users
	res := coll.Find(s.ParentCond(c))
	// TODO: also work with c.Cond for ListCond (limit and offset)
	err = res.All(el)
	return
}

// implements Retrieve method of Service
func (s *Service) Retrieve(c service.Context, el service.EntityListPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	// retrieve all users of id(s)
	res := coll.Find(s.KeyCond(c))
	err = res.All(el)
	return
}

// implements Update method of Service
func (s *Service) Update(c service.Context, e service.EntityPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	res := coll.Find(s.KeyCond(c))

	// update the user
	err = res.Update(e)
	return
}

// implements Delete method of Service
func (s *Service) Delete(c service.Context) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	res := coll.Find(s.KeyCond(c))
	if err != nil {
		return
	}

	err = res.Remove()
	return
}
