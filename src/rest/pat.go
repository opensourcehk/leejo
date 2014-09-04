package rest

import (
	"encoding/json"
	"github.com/gorilla/pat"
	"github.com/gourd/session"
	"leejo/data"
	"log"
	"net/http"
)

// create REST CURD interface with PatRestHelper and pat router
// it knows nothing about the underlying database implementation
// it only handles JSON communication and error handling with http client
func Pat(h Handler, sh session.Handler, p Protocol, r *pat.Router) {

	r.Get(h.EntityPath(), func(w http.ResponseWriter, r *http.Request) {

		var err error

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// allocate memory for variables
		el := s.AllocEntityList()

		// retrieve all of entities in context c
		err = s.Retrieve(c.GetKey(), c.GetParentKey(), el)
		if err != nil {
			Error(w, err)
			return
		}

		// check access
		err = h.CheckAccess("retrieve", sess, el)
		if err != nil {
			Error(w, err)
			return
		}

		if s.Len(el) == 0 {
			w.WriteHeader(404) // not found
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
	r.Get(h.BasePath(), func(w http.ResponseWriter, r *http.Request) {

		var err error

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// allocate memory for variables
		el := s.AllocEntityList()

		// check access
		err = h.CheckAccess("list", sess, c.GetParentKey())
		if err != nil {
			Error(w, err)
			return
		}

		// dummy limit and offset for now
		err = s.List(c.GetKey(), el)
		if err != nil {
			Error(w, err)
			return
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
	r.Post(h.BasePath(), func(w http.ResponseWriter, r *http.Request) {

		// allocate memory for variables
		var err error

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// allocate memory for variables
		e := s.AllocEntity()

		// TODO: find a way to enforce parent key
		err = json.NewDecoder(r.Body).Decode(e)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			Error(w, err)
			return
		}
		log.Printf("Create %#v", e)

		// check access
		err = h.CheckAccess("create", sess, nil)
		if err != nil {
			Error(w, err)
			return
		}

		err = s.Create(c, e)
		if err != nil {
			Error(w, err)
			return
		}

		w.WriteHeader(201) // created
		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Put(h.EntityPath(), func(w http.ResponseWriter, r *http.Request) {

		var err error

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// allocate memory for variables
		e := s.AllocEntity()
		el := s.AllocEntityList()

		err = json.NewDecoder(r.Body).Decode(e)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			Error(w, err)
			return
		}
		log.Printf("Update %#v with %#v", c, e)

		// retrieve all entities with c.Key
		err = s.Retrieve(c.GetKey(), c.GetParentKey(), el)
		if err != nil {
			Error(w, err)
			return
		}

		// check access
		err = h.CheckAccess("update", sess, el)
		if err != nil {
			Error(w, err)
			return
		}

		s.Update(c.GetKey(), c.GetParentKey(), e)

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []interface{}{e},
		})
	})
	r.Delete(h.EntityPath(), func(w http.ResponseWriter, r *http.Request) {

		var err error

		// get service and context
		sess, err := sh.Session(r)
		s := h.Service(sess)
		c := h.Context(sess)

		// allocate memory for variables
		el := s.AllocEntityList()

		// retrieve all entities with c.Key
		err = s.Retrieve(c.GetKey(), c.GetParentKey(), el)
		if err != nil {
			Error(w, err)
			return
		}

		// check access
		err = h.CheckAccess("delete", sess, el)
		if err != nil {
			Error(w, err)
			return
		}

		// delete the item
		err = s.Delete(c.GetKey(), c.GetParentKey())
		if err != nil {
			Error(w, err)
			return
		}

		// remove all results from database
		w.WriteHeader(404) // not found
		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: el,
		})
	})
}
