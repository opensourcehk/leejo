package main

import (
	"flag"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"net/http"
	"regexp"
	"strings"
	"upper.io/db"
	"upper.io/db/postgresql"
)

// The regex to check for the requested format (allows an optional trailing
// slash).
var rxExt = regexp.MustCompile(`(\.(?:xml|text|json))\/?$`)

// Because `panic`s are caught by martini's Recovery handler, it can be used
// to return server-side errors (500). Some helpful text message should probably
// be sent, although not the technical error (which is printed in the log).
func Must(data []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return data
}

func MapEncoder(c martini.Context, w http.ResponseWriter, r *http.Request) {
	// Get the format extension
	matches := rxExt.FindStringSubmatch(r.URL.Path)
	ft := ".json"
	if len(matches) > 1 {
		// Rewrite the URL without the format extension
		l := len(r.URL.Path) - len(matches[1])
		if strings.HasSuffix(r.URL.Path, "/") {
			l--
		}
		r.URL.Path = r.URL.Path[:l]
		ft = matches[1]
	}
	// Inject the requested encoder
	switch ft {
	case ".xml":
		c.MapTo(XmlEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/xml")
	default:
		c.MapTo(JsonEncoder{}, (*Encoder)(nil))
		w.Header().Set("Content-Type", "application/json")
	}
}

func main() {

	// parse flags
	dbhost := flag.String("dbhost", "localhost", "Hostname of the database")
	dbname := flag.String("dbname", "", "Name of the database")
	dbuser := flag.String("dbuser", "", "Name of the database user")
	dbpass := flag.String("dbpass", "", "Password of the database user")
	flag.Parse()

	// connect to database
	var dbsettings = db.Settings{
		Host:     *dbhost,
		Database: *dbname,
		User:     *dbuser,
		Password: *dbpass,
	}
	sess, err := db.Open(postgresql.Adapter, dbsettings)
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	// martini
	m := martini.Classic()
	m.Use(martini.Recovery())

	// render html template from template
	m.Use(render.Renderer())

	// map the MapEncoder middleware
	m.Use(MapEncoder)

	// User
	m.Group("/api.v1/user", func(r martini.Router) {
		r.Get("/:id", func(params martini.Params, enc Encoder, r *http.Request) []byte {
			return Must(enc.Encode("Create User"))
		})
		r.Post("", binding.Bind(User{}), func(user User, enc Encoder) []byte {
			return Must(enc.Encode([]User{user}))
		})
		r.Put("/:id", func(params martini.Params, enc Encoder) []byte {
			return Must(enc.Encode("Replace " + params["id"]))
		})
		r.Delete("/:id", func(params martini.Params, enc Encoder) []byte {
			return Must(enc.Encode("Delete " + params["id"]))
		})
	})

	// Frontend
	m.Group("/", func(r martini.Router) {
		r.Get("", func(r render.Render) {
			r.HTML(200, "home", map[string]interface{}{
				"hello": "world",
			})
		})
	})

	http.Handle("/", m)
	http.ListenAndServe(":8080", nil)

}
