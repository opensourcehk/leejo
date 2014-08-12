package main

import (
	"encoding/json"
	"github.com/RangelReale/osin"
	"github.com/gorilla/pat"
	"io/ioutil"
	"leejo/data"
	"leejo/service"
	"log"
	"net/http"
	"upper.io/db"
)

type UserService struct {
	Db       db.Database
	CollName string
	ApplyKey *func(service.KeyPtr, service.EntityPtr) (err error)
}

func (s *UserService) ApplyCreateKey(k service.KeyPtr, e service.EntityPtr) (err error) {
	u := e.(*data.User)
	u.UserId = k.(int64)
	return
}

func (s *UserService) KeyCond(k service.KeyPtr) db.Cond {
	return db.Cond{"user_id": k}
}

func (s *UserService) ParentCond(pk service.ParentKeyPtr) db.Cond {
	return db.Cond{}
}

func (s *UserService) Create(c service.Context, e service.EntityPtr) (err error) {
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
	err = s.ApplyCreateKey(id, e)
	return
}

func (s *UserService) List(c service.Context, el service.EntityListPtr) (err error) {

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

func (s *UserService) Retrieve(c service.Context, el service.EntityListPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	// retrieve all users of id(s)
	res := coll.Find(s.KeyCond(c.Key))
	err = res.All(el)
	return
}

func (s *UserService) Update(c service.Context, e service.EntityPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	res := coll.Find(s.KeyCond(c.Key))

	// update the user
	err = res.Update(e)
	return
}

func (s *UserService) Delete(c service.Context) (err error) {
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

func bindUser(path string, osinServer *osin.Server, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr

	// return a typical CURD service
	GetService := func(r *http.Request) service.Service {
		return &UserService{
			Db:       sess,
			CollName: "leejo_user",
		}
	}

	// allocate an entity and return the address
	GetEntityPtr := func() service.EntityPtr {
		return &data.User{}
	}

	// allocate a slice of entity and return the address
	GetEntityListPtr := func() service.EntityListPtr {
		return &[]data.User{}
	}

	// translate an http request into a query context
	GetContext := func(r *http.Request) service.Context {
		return service.Context{
			Key:       r.URL.Query().Get(":id"),
			ParentKey: nil,
			Cond: &service.BasicListCond{
				Limit:  20,
				Offset: 0,
			},
		}
	}

	// path which include entity specific information
	subPath := "{id:[0-9]+}"

	r.Get(path+"/"+subPath, func(w http.ResponseWriter, r *http.Request) {
		s := GetService(r)
		el := GetEntityListPtr()
		c := GetContext(r)

		// retrieve all users of c.Key
		err := s.Retrieve(c, el)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		s := GetService(r)
		el := GetEntityListPtr()
		c := GetContext(r)

		// dummy limit and offset for now
		err := s.List(c, el)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
	r.Post(path, func(w http.ResponseWriter, r *http.Request) {
		s := GetService(r)
		e := GetEntityPtr()
		c := GetContext(r)

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, e)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", e)

		err = s.Create(c, e)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Put(path+"/"+subPath, func(w http.ResponseWriter, r *http.Request) {
		s := GetService(r)
		e := GetEntityPtr()
		c := GetContext(r)

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, e)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", e)

		s.Update(c, e)

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Delete(path+"/"+subPath, func(w http.ResponseWriter, r *http.Request) {
		s := GetService(r)
		el := GetEntityListPtr()
		c := GetContext(r)

		// retrieve all entities with c.Key
		err := s.Retrieve(c, el)
		if err != nil {
			panic(err)
		}

		// delete the item
		err = s.Delete(c)
		if err != nil {
			panic(err)
		}

		// remove all results from database
		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
}
