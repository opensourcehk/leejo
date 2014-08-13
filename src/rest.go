package main

import (
	"encoding/json"
	"github.com/gorilla/pat"
	"github.com/gourd/service"
	"io/ioutil"
	"leejo/data"
	"leejo/session"
	"log"
	"net/http"
)

// interface of helper that provides help to create
// a REST CURD interface
type PatRestHelper interface {

	// returns a pat readable regular expression
	// to listing endpoint
	BasePath() string

	// returns a pat readable regular expression
	// to individual entity
	EntityPath() string

	// allocate memory of a single entity and return address
	Entity() service.EntityPtr

	// allocate memory of a slice of entity and return address
	EntityList() service.EntityListPtr

	// allocate storage service for CURD operations
	Service(s session.Session) service.Service

	// translate an http request into a query context
	// i.e. key, parent key, query conditions, limit, offset and etc.
	Context(s session.Session) service.Context

	// check if the session allow
	// the kind of access to this object
	CheckAccess(string, session.Session, interface{}) error
}

// create REST CURD interface with PatRestHelper and pat router
// it knows nothing about the underlying database implementation
// it only handles JSON communication and error handling with http client
func RestOnPat(h PatRestHelper, sh session.SessionHandler, r *pat.Router) {

	r.Get(h.EntityPath(), func(w http.ResponseWriter, r *http.Request) {

		// allocate memory for variables
		var err error
		el := h.EntityList()

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// retrieve all of entities in context c
		err = s.Retrieve(c.GetKey(), c.GetParentKey(), el)
		if err != nil {
			panic(err)
		}

		// check access
		err = h.CheckAccess("retrieve", sess, el)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
	r.Get(h.BasePath(), func(w http.ResponseWriter, r *http.Request) {

		// allocate memory for variables
		var err error
		el := h.EntityList()

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// check access
		err = h.CheckAccess("list", sess, c.GetParentKey())
		if err != nil {
			panic(err)
		}

		// dummy limit and offset for now
		err = s.List(c.GetKey(), el)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
	r.Post(h.BasePath(), func(w http.ResponseWriter, r *http.Request) {

		// allocate memory for variables
		var err error
		e := h.Entity()

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// TODO: find a way to enforce parent key

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

		// check access
		err = h.CheckAccess("create", sess, nil)
		if err != nil {
			panic(err)
		}

		err = s.Create(c, e)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Put(h.EntityPath(), func(w http.ResponseWriter, r *http.Request) {

		// allocate memory for variables
		var err error
		e := h.Entity()
		el := h.EntityList()

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

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

		// retrieve all entities with c.Key
		err = s.Retrieve(c.GetKey(), c.GetParentKey(), el)
		if err != nil {
			panic(err)
		}

		// check access
		err = h.CheckAccess("update", sess, el)
		if err != nil {
			panic(err)
		}

		s.Update(c.GetKey(), c.GetParentKey(), e)

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Delete(h.EntityPath(), func(w http.ResponseWriter, r *http.Request) {

		// allocate memory for variables
		var err error
		el := h.EntityList()

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// retrieve all entities with c.Key
		err = s.Retrieve(c.GetKey(), c.GetParentKey(), el)
		if err != nil {
			panic(err)
		}

		// check access
		err = h.CheckAccess("delete", sess, el)
		if err != nil {
			panic(err)
		}

		// delete the item
		err = s.Delete(c.GetKey(), c.GetParentKey())
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
