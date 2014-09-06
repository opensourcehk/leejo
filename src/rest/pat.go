package rest

import (
	"github.com/gorilla/pat"
	"github.com/gourd/session"
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
			WriteError(w, sess, p, err)
			return
		}

		// check access
		err = h.CheckAccess("retrieve", sess, el)
		if err != nil {
			WriteError(w, sess, p, err)
			return
		}

		// determine the status by query result
		status := http.StatusOK
		if s.Len(el) == 0 {
			status = http.StatusNotFound
		}

		// write response with proper status code
		WriteResponse(w, sess, p, el, status)
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
			WriteError(w, sess, p, err)
			return
		}

		// dummy limit and offset for now
		err = s.List(c.GetKey(), el)
		if err != nil {
			WriteError(w, sess, p, err)
			return
		}

		// write response with status "ok"
		WriteResponse(w, sess, p, el, http.StatusOK)
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
		err = p.NewDecoder(sess, r.Body).Decode(e)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			WriteError(w, sess, p, err)
			return
		}
		log.Printf("Create %#v", e)

		// check access
		err = h.CheckAccess("create", sess, nil)
		if err != nil {
			WriteError(w, sess, p, err)
			return
		}

		// create the entity
		err = s.Create(c, e)
		if err != nil {
			WriteError(w, sess, p, err)
			return
		}

		// write response with status "created"
		WriteResponse(w, sess, p, []interface{}{e}, http.StatusCreated)
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

		err = p.NewDecoder(sess, r.Body).Decode(e)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			WriteError(w, sess, p, err)
			return
		}
		log.Printf("Update %#v with %#v", c, e)

		// retrieve all entities with c.Key
		err = s.Retrieve(c.GetKey(), c.GetParentKey(), el)
		if err != nil {
			WriteError(w, sess, p, err)
			return
		}

		// check access
		err = h.CheckAccess("update", sess, el)
		if err != nil {
			WriteError(w, sess, p, err)
			return
		}

		// update the entity
		s.Update(c.GetKey(), c.GetParentKey(), e)

		// write response with status "ok"
		WriteResponse(w, sess, p, []interface{}{e}, http.StatusOK)
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
			WriteError(w, sess, p, err)
			return
		}

		// check access
		err = h.CheckAccess("delete", sess, el)
		if err != nil {
			WriteError(w, sess, p, err)
			return
		}

		// delete the item
		err = s.Delete(c.GetKey(), c.GetParentKey())
		if err != nil {
			WriteError(w, sess, p, err)
			return
		}

		// remove all results from database
		// then report "not found" to client
		WriteResponse(w, sess, p, el, http.StatusNotFound)
	})
}
