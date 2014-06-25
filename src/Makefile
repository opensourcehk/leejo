build: martini martini-binding martini-render upper-db upper-db-pgsql
	go fmt
	go build -o ${BIN}/leejo_server

martini: github.com/go-martini/martini

martini-binding: github.com/martini-contrib/binding

martini-render: github.com/martini-contrib/render

gopostgresql: github.com/xiam/gopostgresql

gosexy-to: menteslibres.net/gosexy/to

upper-db: upper.io/db gosexy-to

upper-db-pgsql: upper.io/db/postgresql gopostgresql

github.com/go-martini/martini:
	go get github.com/go-martini/martini

github.com/martini-contrib/binding:
	go get github.com/martini-contrib/binding

github.com/martini-contrib/render:
	go get github.com/martini-contrib/render

github.com/xiam/gopostgresql:
	go get github.com/xiam/gopostgresql

menteslibres.net/gosexy/to:
	go get menteslibres.net/gosexy/to

upper.io/db:
	go get upper.io/db

upper.io/db/postgresql:
	go get upper.io/postgresql

.PHONY: build martini martini-binding martini-render upper-db upper-db-pgsql gopostgresql gosexy-to
