package main

import (
	"data"
	"encoding/json"
	"github.com/RangelReale/osin"
	"github.com/gorilla/pat"
	"io/ioutil"
	"log"
	"net/http"
	"upper.io/db"
)

func bindUser(path string, osinServer *osin.Server, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		userCollection, err := sess.Collection("leejo_user")
		if err != nil {
			panic(err)
		}

		// retrieve all users
		res := userCollection.Find()
		var users []data.User
		err = res.All(&users)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: users,
		})
	})
	r.Get(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		userCollection, err := sess.Collection("leejo_user")
		if err != nil {
			panic(err)
		}

		// retrieve all users of id(s)
		id := r.URL.Query().Get(":id")
		res := userCollection.Find(db.Cond{"user_id": id})
		var users []data.User
		err = res.All(&users)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: users,
		})
	})
	r.Post(path, func(w http.ResponseWriter, r *http.Request) {

		user := data.User{}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, &user)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", user)

		userCollection, err := sess.Collection("leejo_user")
		if err != nil {
			panic(err)
		}

		// add the user to user collection
		userId, err := userCollection.Append(user)
		if err != nil {
			panic(err)
		}
		user.UserId = userId.(int64)

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []data.User{user},
		})
	})
	r.Put(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {

		user := data.User{}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, &user)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", user)

		var users []data.User
		userCollection, err := sess.Collection("leejo_user")
		if err != nil {
			panic(err)
		}

		id := r.URL.Query().Get(":id")
		res := userCollection.Find(db.Cond{"user_id": id})

		// update the user
		err = res.Update(user)
		if err != nil {
			panic(err)
		}

		// retrieve the just updated record from database
		res = userCollection.Find(db.Cond{"user_id": id})
		err = res.All(&users)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: users,
		})
	})
	r.Delete(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		userCollection, err := sess.Collection("leejo_user")
		if err != nil {
			panic(err)
		}

		// retrieve all users of id(s)
		id := r.URL.Query().Get(":id")
		res := userCollection.Find(db.Cond{"user_id": id})
		var users []data.User
		err = res.All(&users)
		if err != nil {
			panic(err)
		}

		// TODO: if len(users) == 0, return 404 error

		// remove all results from database
		err = res.Remove()
		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: users,
		})

	})
}
