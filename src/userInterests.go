package main

import (
	"data"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"net/http"
	"strconv"
	"upper.io/db"
)

func bindUserInterests(path string, sessPtr *db.Database, m *martini.ClassicMartini) {
	sess := *sessPtr
	m.Group(path, func(r martini.Router) {
		r.Get("", func(params martini.Params, enc Encoder, r *http.Request) []byte {
			userInterestsCol, err := sess.Collection("leejo_user_interest")
			if err != nil {
				panic(err)
			}

			// retrieve all userInterests of the user_id
			res := userInterestsCol.Find(db.Cond{
				"user_id": params["user_id"],
			})
			var userInterests []data.UserInterest
			err = res.All(&userInterests)
			if err != nil {
				panic(err)
			}

			return Must(enc.Encode(Resp{
				Status: "OK",
				Result: userInterests,
			}))
		})
		r.Get("/:id", func(params martini.Params, enc Encoder, r *http.Request) []byte {
			userInterestsCol, err := sess.Collection("leejo_user_interest")
			if err != nil {
				panic(err)
			}

			// retrieve all userInterests of the user_id and id(s)
			res := userInterestsCol.Find(db.Cond{
				"user_interest_id": params["id"],
				"user_id":          params["user_id"],
			})
			var userInterests []data.UserInterest
			err = res.All(&userInterests)
			if err != nil {
				panic(err)
			}

			return Must(enc.Encode(Resp{
				Status: "OK",
				Result: userInterests,
			}))
		})
		r.Post("", binding.Bind(data.UserInterest{}), func(
			params martini.Params, user data.UserInterest, enc Encoder) []byte {

			inputUserId, err := strconv.ParseInt(params["user_id"], 10, 64)
			if err != nil {
				panic(err)
			}
			user.UserId = inputUserId // force to use URL's user id

			userInterestsCol, err := sess.Collection("leejo_user_interest")
			if err != nil {
				panic(err)
			}

			// add the user to user collection
			userId, err := userInterestsCol.Append(user)
			if err != nil {
				panic(err)
			}
			user.UserInterestId = userId.(int64)

			return Must(enc.Encode(Resp{
				Status: "OK",
				Result: []data.UserInterest{user},
			}))
		})
		r.Put("/:id", binding.Bind(data.UserInterest{}), func(user data.UserInterest, params martini.Params, enc Encoder) []byte {

			var userInterests []data.UserInterest
			userInterestsCol, err := sess.Collection("leejo_user_interest")
			if err != nil {
				panic(err)
			}

			res := userInterestsCol.Find(db.Cond{
				"user_interest_id": params["id"],
				"user_id":          params["user_id"],
			})

			// update the user
			err = res.Update(user)
			if err != nil {
				panic(err)
			}

			// retrieve the just updated record from database
			res = userInterestsCol.Find(db.Cond{
				"user_interest_id": params["id"],
				"user_id":          params["user_id"],
			})
			err = res.All(&userInterests)
			if err != nil {
				panic(err)
			}

			return Must(enc.Encode(Resp{
				Status: "OK",
				Result: userInterests,
			}))
		})
		r.Delete("/:id", func(params martini.Params, enc Encoder) []byte {
			userInterestsCol, err := sess.Collection("leejo_user_interest")
			if err != nil {
				panic(err)
			}

			// retrieve all userInterests of id(s)
			res := userInterestsCol.Find(db.Cond{
				"user_interest_id": params["id"],
				"user_id":          params["user_id"],
			})
			var userInterests []data.UserInterest
			err = res.All(&userInterests)
			if err != nil {
				panic(err)
			}

			// TODO: if len(userInterests) == 0, return 404 error

			// remove all results from database
			err = res.Remove()
			return Must(enc.Encode(Resp{
				Status: "OK",
				Result: userInterests,
			}))

		})
	})
}
