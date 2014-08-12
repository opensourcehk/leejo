export ROOT=$(shell pwd)
export GOPATH=${ROOT}/gopath
export BIN=${ROOT}/bin


#
# main targets
#

all: build

serve: build
	./bin/leejo_server -config ./data/config.json

build: bin/leejo_server

test: bin/integration_test
	@echo "Integration Test"
	@bin/integration_test

check:
	@echo "Unit Test"
	@cd src; go test
	@cd src/data; go test
	@cd src/oauth2; go test
	@cd tests; go test

clean:
	rm -Rf bin/*

.PHONY: all serve build check test clean



#
# server
#

bin/leejo_server: pat gourd-service osin upper-db-pgsql
	cd src; go build -o ${BIN}/leejo_server

pat: gopath/src/github.com/gorilla/pat

gourd-service: gopath/src/github.com/gourd/service

osin: gopath/src/github.com/RangelReale/osin

upper-db-pgsql: gopath/src/upper.io/db gopath/src/menteslibres.net/gosexy/to gopath/src/upper.io/db/postgresql gopath/src/github.com/xiam/gopostgresql

gopath/src/github.com/gorilla/pat:
	go get github.com/gorilla/pat

gopath/src/github.com/gourd/service:
	go get github.com/gourd/service

gopath/src/github.com/go-martini/martini:
	go get github.com/go-martini/martini

gopath/src/github.com/martini-contrib/binding:
	go get github.com/martini-contrib/binding

gopath/src/github.com/martini-contrib/render:
	go get github.com/martini-contrib/render

gopath/src/github.com/xiam/gopostgresql:
	go get github.com/xiam/gopostgresql

gopath/src/menteslibres.net/gosexy/to:
	go get menteslibres.net/gosexy/to

gopath/src/upper.io/db:
	go get upper.io/db

gopath/src/upper.io/db/postgresql:
	go get upper.io/db/postgresql

gopath/src/github.com/RangelReale/osin:
	go get github.com/RangelReale/osin

.PHONY: martini osin upper-db-pgsql


#
# tests
#

bin/integration_test: \
	gopath/src/github.com/jmcvetta/napping \
	gopath/src/github.com/yookoala/restit \
	gopath/src/github.com/skratchdot/open-golang/open
	cd tests; go build -o ${BIN}/integration_test

gopath/src/github.com/yookoala/restit:
	go get github.com/yookoala/restit

gopath/src/github.com/jmcvetta/napping:
	go get github.com/jmcvetta/napping

gopath/src/github.com/skratchdot/open-golang/open:
	go get github.com/skratchdot/open-golang/open
