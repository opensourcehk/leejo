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

func (s *UserService) Create(e service.EntityPtr) (err error) {
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

func (s *UserService) List(pk service.ParentKeyPtr,
	c service.ListCond, el service.EntityListPtr) (err error) {

	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	// retrieve all users
	res := coll.Find(s.ParentCond(pk))
	err = res.All(el)
	return
}

func (s *UserService) Retrieve(k service.KeyPtr, el service.EntityListPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	// retrieve all users of id(s)
	res := coll.Find(s.KeyCond(k))
	err = res.All(el)
	return
}

func (s *UserService) Update(k service.KeyPtr, e service.EntityPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	res := coll.Find(s.KeyCond(k))

	// update the user
	err = res.Update(e)
	return
}

func (s *UserService) Delete(k service.KeyPtr) (err error) {
	coll, err := s.Db.Collection(s.CollName)
	if err != nil {
		return
	}

	res := coll.Find(s.KeyCond(k))
	if err != nil {
		return
	}

	err = res.Remove()
	return
}

func bindUser(path string, osinServer *osin.Server, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr

	GetUserService := func(r *http.Request) service.Service {
		return &UserService{
			Db: sess,
			CollName: "leejo_user",
		}
	}

	r.Get(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		s := GetUserService(r)
		el := []data.User{}

		// retrieve all users of id(s)
		id := r.URL.Query().Get(":id")
		err := s.Retrieve(id, &el)
		s.Retrieve(id, &el)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		s := GetUserService(r)
		el := []data.User{}

		// dummy limit and offset for now
		c := service.BasicListCond{
			Limit: 20,
		}
		err := s.List(nil, &c, &el)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
	r.Post(path, func(w http.ResponseWriter, r *http.Request) {
		s := GetUserService(r)
		e := data.User{}

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, &e)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", e)

		err = s.Create(&e)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []data.User{e},
		})
	})
	r.Put(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		s := GetUserService(r)
		e := data.User{}

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, &e)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", e)

		id := r.URL.Query().Get(":id")
		s.Update(id, &e)

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Delete(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		s := GetUserService(r)
		el := []data.User{}

		// retrieve all users of id(s)
		id := r.URL.Query().Get(":id")
		err := s.Retrieve(id, &el)
		if err != nil {
			panic(err)
		}

		// delete the item
		err = s.Delete(id)
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
