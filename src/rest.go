package main

import (
	"encoding/json"
	"github.com/gorilla/pat"
	"io/ioutil"
	"leejo/data"
	"leejo/service"
	"log"
	"net/http"
)

// interface of helper that provides help to create
// a REST CURD interface
type PatRestHelper interface {

	// returns a pat readable regular expression
	// of an individual entity
	SubPath() string

	// allocate storage service for CURD operations
	Service(r *http.Request) service.Service

	// translate an http request into a query context
	// i.e. key, parent key, query conditions, limit, offset and etc.
	Context(r *http.Request) service.Context

	// allocate memory of a single entity and return address
	Entity() service.EntityPtr

	// allocate memory of a slice of entity and return address
	EntityList() service.EntityListPtr
}

// create REST CURD interface with PatRestHelper and pat router
// it knows nothing about the underlying database implementation
// it only handles JSON communication and error handling with http client
func RestOnPat(path string, h PatRestHelper, r *pat.Router) {

	r.Get(path+"/"+h.SubPath(), func(w http.ResponseWriter, r *http.Request) {
		s := h.Service(r)
		el := h.EntityList()
		c := h.Context(r)

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
		s := h.Service(r)
		el := h.EntityList()
		c := h.Context(r)

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
		s := h.Service(r)
		e := h.Entity()
		c := h.Context(r)

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
		log.Printf("Create %#v", e)

		err = s.Create(c, e)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Put(path+"/"+h.SubPath(), func(w http.ResponseWriter, r *http.Request) {
		s := h.Service(r)
		e := h.Entity()
		c := h.Context(r)

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
		log.Printf("Update %#v with %#v", c, e)

		s.Update(c, e)

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Delete(path+"/"+h.SubPath(), func(w http.ResponseWriter, r *http.Request) {
		s := h.Service(r)
		el := h.EntityList()
		c := h.Context(r)

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
