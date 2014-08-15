export ROOT=$(shell pwd)
export GOPATH=${ROOT}/gopath
export BIN=${ROOT}/bin


#
# main targets
#

all: check build

serve: build
	@echo "Serve"
	@echo "====="
	./bin/leejo_server -config ./data/config.json

build: fmt bin/leejo_server bin/integration_test

test: fmt bin/integration_test
	@echo "Integration Test"
	@echo "================"
	@bin/integration_test
	@echo

update: clean-gopath get-deps get-test-deps

clean-gopath:
	rm -Rf gopath/src/github.com
	rm -Rf gopath/src/code.google.com
	rm -Rf gopath/src/menteslibres.net
	rm -Rf gopath/src/upper.io

check: get-deps get-test-deps
	@echo "Unit Test"
	@echo "========="
	@cd src; go test -i; go test
	@cd src/data; go test -i; go test
	@cd src/oauth2; go test -i; go test
	@cd src/session; go test -i; go test
	@cd tests; go test -i; go test
	@echo

fmt:
	@echo "Format Code"
	@echo "==========="
	cd src; go fmt
	cd src/data; go fmt
	cd src/oauth2; go fmt
	cd src/session; go fmt
	cd tests; go fmt
	@echo

clean:
	rm -Rf bin/*

.PHONY: all update clean-gopath serve build test check fmt clean



#
# server
#

bin/leejo_server: get-deps
	@echo "Build Server"
	@echo "============"
	cd src; go test -i
	cd src; go build -o ${BIN}/leejo_server
	@echo

get-deps: pat gourd-service osin upper-db-pgsql

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

.PHONY: pat gourd-service osin upper-db-pgsql


#
# tests
#

get-test-deps: \
	gopath/src/github.com/jmcvetta/napping \
	gopath/src/github.com/yookoala/restit \
	gopath/src/github.com/skratchdot/open-golang/open

bin/integration_test: get-test-deps
	@echo "Build Integration Test"
	@echo "======================"
	cd tests; go build -o ${BIN}/integration_test
	@echo

gopath/src/github.com/yookoala/restit:
	go get github.com/yookoala/restit

gopath/src/github.com/jmcvetta/napping:
	go get github.com/jmcvetta/napping

gopath/src/github.com/skratchdot/open-golang/open:
	go get github.com/skratchdot/open-golang/open

.PHONY: get-test-deps bin/integration_test
