package main

import (
	"data"
	"github.com/RangelReale/osin"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"net/http"
	"upper.io/db"
)

func bindUser(path string, osinServer *osin.Server, sessPtr *db.Database, m *martini.ClassicMartini) {
	sess := *sessPtr
	m.Group(path, func(r martini.Router) {
		r.Get("", func(params martini.Params, enc Encoder, r *http.Request) []byte {
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

			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: users,
			}))
		})
		r.Get("/:id", func(params martini.Params, enc Encoder, r *http.Request) []byte {
			userCollection, err := sess.Collection("leejo_user")
			if err != nil {
				panic(err)
			}

			// retrieve all users of id(s)
			res := userCollection.Find(db.Cond{"user_id": params["id"]})
			var users []data.User
			err = res.All(&users)
			if err != nil {
				panic(err)
			}

			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: users,
			}))
		})
		r.Post("", binding.Bind(data.User{}), func(user data.User, enc Encoder) []byte {
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

			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: []data.User{user},
			}))
		})
		r.Put("/:id", binding.Bind(data.User{}), func(user data.User, params martini.Params, enc Encoder) []byte {

			var users []data.User
			userCollection, err := sess.Collection("leejo_user")
			if err != nil {
				panic(err)
			}

			res := userCollection.Find(db.Cond{"user_id": params["id"]})

			// update the user
			err = res.Update(user)
			if err != nil {
				panic(err)
			}

			// retrieve the just updated record from database
			res = userCollection.Find(db.Cond{"user_id": params["id"]})
			err = res.All(&users)
			if err != nil {
				panic(err)
			}

			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: users,
			}))
		})
		r.Delete("/:id", func(params martini.Params, enc Encoder) []byte {
			userCollection, err := sess.Collection("leejo_user")
			if err != nil {
				panic(err)
			}

			// retrieve all users of id(s)
			res := userCollection.Find(db.Cond{"user_id": params["id"]})
			var users []data.User
			err = res.All(&users)
			if err != nil {
				panic(err)
			}

			// TODO: if len(users) == 0, return 404 error

			// remove all results from database
			err = res.Remove()
			return Must(enc.Encode(data.Resp{
				Status: "OK",
				Result: users,
			}))

		})
	})
}
