package main

import (
	"encoding/json"
	"github.com/RangelReale/osin"
	"github.com/gorilla/pat"
	"io/ioutil"
	"leejo/data"
	"log"
	"net/http"
	"strconv"
	"upper.io/db"
)

func bindUserInterests(path string, osinServer *osin.Server, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		userInterestsCol, err := sess.Collection("leejo_user_interest")
		if err != nil {
			panic(err)
		}

		// retrieve all userInterests of the user_id
		user_id := r.URL.Query().Get(":user_id")
		res := userInterestsCol.Find(db.Cond{
			"user_id": user_id,
		})
		var userInterests []data.UserInterest
		err = res.All(&userInterests)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: userInterests,
		})
	})
	r.Get(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		userInterestsCol, err := sess.Collection("leejo_user_interest")
		if err != nil {
			panic(err)
		}

		// retrieve all userInterests of the user_id and id(s)
		id := r.URL.Query().Get(":id")
		user_id := r.URL.Query().Get(":user_id")
		res := userInterestsCol.Find(db.Cond{
			"user_interest_id": id,
			"user_id":          user_id,
		})
		var userInterests []data.UserInterest
		err = res.All(&userInterests)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: userInterests,
		})
	})
	r.Post(path, func(w http.ResponseWriter, r *http.Request) {

		userInterest := data.UserInterest{}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, &userInterest)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", userInterest)

		user_id := r.URL.Query().Get(":user_id")
		inputUserId, err := strconv.ParseInt(user_id, 10, 64)
		if err != nil {
			panic(err)
		}
		userInterest.UserId = inputUserId // force to use URL's user id

		userInterestsCol, err := sess.Collection("leejo_user_interest")
		if err != nil {
			panic(err)
		}

		// add the user to user collection
		userId, err := userInterestsCol.Append(userInterest)
		if err != nil {
			panic(err)
		}
		userInterest.UserInterestId = userId.(int64)

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []data.UserInterest{userInterest},
		})
	})
	r.Put(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {

		userInterest := data.UserInterest{}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, &userInterest)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", userInterest)

		var userInterests []data.UserInterest
		userInterestsCol, err := sess.Collection("leejo_user_interest")
		if err != nil {
			panic(err)
		}

		id := r.URL.Query().Get(":id")
		user_id := r.URL.Query().Get(":user_id")
		res := userInterestsCol.Find(db.Cond{
			"user_interest_id": id,
			"user_id":          user_id,
		})

		// update the user
		err = res.Update(userInterest)
		if err != nil {
			panic(err)
		}

		// retrieve the just updated record from database
		res = userInterestsCol.Find(db.Cond{
			"user_interest_id": id,
			"user_id":          user_id,
		})
		err = res.All(&userInterests)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: userInterests,
		})
	})
	r.Delete(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		userInterestsCol, err := sess.Collection("leejo_user_interest")
		if err != nil {
			panic(err)
		}

		// retrieve all userInterests of id(s)
		id := r.URL.Query().Get(":id")
		user_id := r.URL.Query().Get(":user_id")
		res := userInterestsCol.Find(db.Cond{
			"user_interest_id": id,
			"user_id":          user_id,
		})
		var userInterests []data.UserInterest
		err = res.All(&userInterests)
		if err != nil {
			panic(err)
		}

		// TODO: if len(userInterests) == 0, return 404 error

		// remove all results from database
		err = res.Remove()
		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: userInterests,
		})

	})
}
