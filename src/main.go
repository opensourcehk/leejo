package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"io/ioutil"
	"net/http"
	"os"
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

type config struct {
	Db dbconfig `json:"db"`
}

type dbconfig struct {
	Host string `json:"host"`
	Name string `json:"name"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func main() {

	// parse flags
	confFn := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	// read the config file to conf
	confFile, err := ioutil.ReadFile(*confFn)
	if err != nil {
		fmt.Printf("Failed opening config file \"%s\": %v\n", *confFn, err)
		os.Exit(1)
	}
	conf := config{}
	json.Unmarshal(confFile, &conf)

	// connect to database
	var dbsettings = db.Settings{
		Host:     conf.Db.Host,
		Database: conf.Db.Name,
		User:     conf.Db.User,
		Password: conf.Db.Pass,
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

	// Users related API
	bindUser("/api.v1/users", &sess, m)

	// UserSkills related API
	bindUserSkills("/api.v1/userSkills/:user_id", &sess, m)

	// UserInterests related API
	bindUserInterests("/api.v1/userInterests/:user_id", &sess, m)

	// handle OAuth2 endpoints
	bindAuth("/oauth2", &sess)

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
