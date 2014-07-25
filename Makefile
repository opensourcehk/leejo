
GOPATH=$(shell pwd)
BIN=${GOPATH}/bin

all: build

run: build
	./bin/leejo_server -config ./data/config.json

build: bin/leejo_server

check: bin/integration_test
	@echo "Integration Test"
	@bin/integration_test

test:
	@echo "Unit Test"
	cd src; make GOPATH="${GOPATH}" BIN="${BIN}" test

bin/leejo_server:
	cd src; make GOPATH="${GOPATH}" BIN="${BIN}" build

bin/integration_test:
	cd tests; make GOPATH="${GOPATH}" BIN="${BIN}" build

clean:
	rm -Rf bin/*

.PHONY: all run build check test clean
