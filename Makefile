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
	rm -Rf gopath/pkg
	rm -Rf gopath/src/github.com
	rm -Rf gopath/src/code.google.com
	rm -Rf gopath/src/menteslibres.net
	rm -Rf gopath/src/upper.io

check: get-deps get-test-deps
	@echo "Unit Test"
	@echo "========="
	@cd src; go test -i; go test
	@cd lib/data; go test -i; go test
	@cd lib/oauth2; go test -i; go test
	@cd lib/rest; go test -i; go test
	@cd tests; go test -i; go test
	@echo

fmt:
	@echo "Format Code"
	@echo "==========="
	cd src; go fmt
	cd lib/data; go fmt
	cd lib/oauth2; go fmt
	cd tests; go fmt
	@echo

clean:
	rm -Rf gopath/pkg
	rm -Rf gopath/src/github.com/opensourcehk/leejo
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

get-deps: pat gourd osin upper-db-pgsql leejo
	@echo

leejo: gopath/src/github.com/opensourcehk/leejo
	go install github.com/opensourcehk/leejo/lib/data
	go install github.com/opensourcehk/leejo/lib/oauth2
	go install github.com/opensourcehk/leejo/lib/rest

pat: gopath/src/github.com/gorilla/pat

gourd: \
	gopath/src/github.com/gourd/service \
	gopath/src/github.com/gourd/session

osin: gopath/src/github.com/RangelReale/osin

upper-db-pgsql: \
	gopath/src/menteslibres.net/gosexy/to \
	gopath/src/upper.io/db/postgresql \
	gopath/src/github.com/xiam/gopostgresql

gopath/src/github.com/opensourcehk/leejo:
	mkdir -p gopath/src/github.com/opensourcehk/leejo
	ln -s ../../../../../lib gopath/src/github.com/opensourcehk/leejo/lib

gopath/src/github.com/gorilla/pat:
	go get -u github.com/gorilla/pat

gopath/src/github.com/gourd/service:
	go get -u github.com/gourd/service

gopath/src/github.com/gourd/session:
	go get -u github.com/gourd/session

gopath/src/github.com/go-martini/martini:
	go get -u github.com/go-martini/martini

gopath/src/github.com/martini-contrib/binding:
	go get -u github.com/martini-contrib/binding

gopath/src/github.com/martini-contrib/render:
	go get -u github.com/martini-contrib/render

gopath/src/github.com/xiam/gopostgresql:
	go get -u github.com/xiam/gopostgresql

gopath/src/menteslibres.net/gosexy/to:
	go get -u menteslibres.net/gosexy/to

gopath/src/upper.io/db/postgresql:
	go get -u code.google.com/p/go-uuid/uuid
	go get -u upper.io/cache
	go get -u upper.io/db
	go get -u upper.io/db/postgresql

gopath/src/github.com/RangelReale/osin:
	go get -u github.com/RangelReale/osin

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
