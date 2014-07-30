package main

import (
	"data"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"upper.io/db"
)

func bindUserSkills(path string, sessPtr *db.Database, m *martini.ClassicMartini) {
	sess := *sessPtr
	m.Group(path, func(r martini.Router) {
		r.Get("", func(params martini.Params, enc Encoder, r *http.Request) []byte {
			userSkillsCol, err := sess.Collection("leejo_user_skill")
			if err != nil {
				panic(err)
			}

			// retrieve all userSkills of the user_id
			res := userSkillsCol.Find(db.Cond{
				"user_id": params["user_id"],
			})
			var userSkills []data.UserSkill
			err = res.All(&userSkills)
			if err != nil {
				panic(err)
			}

			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: userSkills,
			}))
		})
		r.Get("/:id", func(params martini.Params, enc Encoder, r *http.Request) []byte {
			userSkillsCol, err := sess.Collection("leejo_user_skill")
			if err != nil {
				panic(err)
			}

			// retrieve all userSkills of the user_id and id(s)
			res := userSkillsCol.Find(db.Cond{
				"user_skill_id": params["id"],
				"user_id":       params["user_id"],
			})
			var userSkills []data.UserSkill
			err = res.All(&userSkills)
			if err != nil {
				panic(err)
			}

			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: userSkills,
			}))
		})
		r.Post("", binding.Bind(data.UserSkill{}), func(
			params martini.Params, user data.UserSkill, enc Encoder) []byte {

			inputUserId, err := strconv.ParseInt(params["user_id"], 10, 64)
			if err != nil {
				panic(err)
			}
			user.UserId = inputUserId // force to use URL's user id

			userSkillsCol, err := sess.Collection("leejo_user_skill")
			if err != nil {
				panic(err)
			}

			// add the user to user collection
			userId, err := userSkillsCol.Append(user)
			if err != nil {
				panic(err)
			}
			user.UserSkillId = userId.(int64)

			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: []data.UserSkill{user},
			}))
		})
		r.Put("/:id", binding.Bind(data.UserSkill{}), func(user data.UserSkill, params martini.Params, enc Encoder) []byte {

			var userSkills []data.UserSkill
			userSkillsCol, err := sess.Collection("leejo_user_skill")
			if err != nil {
				panic(err)
			}

			res := userSkillsCol.Find(db.Cond{
				"user_skill_id": params["id"],
				"user_id":       params["user_id"],
			})

			// update the user
			err = res.Update(user)
			if err != nil {
				panic(err)
			}

			// retrieve the just updated record from database
			res = userSkillsCol.Find(db.Cond{
				"user_skill_id": params["id"],
				"user_id":       params["user_id"],
			})
			err = res.All(&userSkills)
			if err != nil {
				panic(err)
			}

			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: userSkills,
			}))
		})
		r.Delete("/:id", func(params martini.Params, enc Encoder) []byte {
			userSkillsCol, err := sess.Collection("leejo_user_skill")
			if err != nil {
				panic(err)
			}

			// retrieve all userSkills of id(s)
			res := userSkillsCol.Find(db.Cond{
				"user_skill_id": params["id"],
				"user_id":       params["user_id"],
			})
			var userSkills []data.UserSkill
			err = res.All(&userSkills)
			if err != nil {
				panic(err)
			}

			// TODO: if len(userSkills) == 0, return 404 error

			// remove all results from database
			err = res.Remove()
			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: userSkills,
			}))

		})
	})
}
