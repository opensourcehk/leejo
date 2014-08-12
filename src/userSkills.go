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

func bindUserSkills(path string, osinServer *osin.Server, sessPtr *db.Database, r *pat.Router) {
	sess := *sessPtr
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		userSkillsCol, err := sess.Collection("leejo_user_skill")
		if err != nil {
			panic(err)
		}

		// retrieve all userSkills of the user_id
		user_id := r.URL.Query().Get(":user_id")
		res := userSkillsCol.Find(db.Cond{
			"user_id": user_id,
		})
		var userSkills []data.UserSkill
		err = res.All(&userSkills)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: userSkills,
		})
	})
	r.Get(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		userSkillsCol, err := sess.Collection("leejo_user_skill")
		if err != nil {
			panic(err)
		}

		// retrieve all userSkills of the user_id and id(s)
		id := r.URL.Query().Get(":id")
		user_id := r.URL.Query().Get(":user_id")
		res := userSkillsCol.Find(db.Cond{
			"user_skill_id": id,
			"user_id":       user_id,
		})
		var userSkills []data.UserSkill
		err = res.All(&userSkills)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: userSkills,
		})
	})
	r.Post(path, func(w http.ResponseWriter, r *http.Request) {

		userSkill := data.UserSkill{}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, &userSkill)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", userSkill)

		user_id := r.URL.Query().Get(":user_id")
		inputUserId, err := strconv.ParseInt(user_id, 10, 64)
		if err != nil {
			panic(err)
		}
		userSkill.UserId = inputUserId // force to use URL's user id

		userSkillsCol, err := sess.Collection("leejo_user_skill")
		if err != nil {
			panic(err)
		}

		// add the user to user collection
		userId, err := userSkillsCol.Append(userSkill)
		if err != nil {
			panic(err)
		}
		userSkill.UserSkillId = userId.(int64)

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: []data.UserSkill{userSkill},
		})
	})
	r.Put(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {

		userSkill := data.UserSkill{}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request: ", err)
			w.WriteHeader(500)
			return
		}

		err = json.Unmarshal(bytes, &userSkill)
		if err != nil {
			log.Printf("Error JSON Unmarshal: ", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("Request %#v", userSkill)

		var userSkills []data.UserSkill
		userSkillsCol, err := sess.Collection("leejo_user_skill")
		if err != nil {
			panic(err)
		}

		id := r.URL.Query().Get(":id")
		user_id := r.URL.Query().Get(":user_id")
		res := userSkillsCol.Find(db.Cond{
			"user_skill_id": id,
			"user_id":       user_id,
		})

		// update the user
		err = res.Update(userSkill)
		if err != nil {
			panic(err)
		}

		// retrieve the just updated record from database
		res = userSkillsCol.Find(db.Cond{
			"user_skill_id": id,
			"user_id":       user_id,
		})
		err = res.All(&userSkills)
		if err != nil {
			panic(err)
		}

		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: userSkills,
		})
	})
	r.Delete(path+"/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		userSkillsCol, err := sess.Collection("leejo_user_skill")
		if err != nil {
			panic(err)
		}

		// retrieve all userSkills of id(s)
		id := r.URL.Query().Get(":id")
		user_id := r.URL.Query().Get(":user_id")
		res := userSkillsCol.Find(db.Cond{
			"user_skill_id": id,
			"user_id":       user_id,
		})
		var userSkills []data.UserSkill
		err = res.All(&userSkills)
		if err != nil {
			panic(err)
		}

		// TODO: if len(userSkills) == 0, return 404 error

		// remove all results from database
		err = res.Remove()
		json.NewEncoder(w).Encode(data.Resp{
			Status: "OK",
			Result: userSkills,
		})

	})
}
