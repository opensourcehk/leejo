
GOPATH=$(shell pwd)
BIN=${GOPATH}/bin

all: build

run: build

build: bin/leejo_server

check: bin/integration_test
	bin/integration_test

bin/leejo_server:
	cd src; make GOPATH="${GOPATH}" BIN="${BIN}" build

bin/integration_test:
	cd tests; make GOPATH="${GOPATH}" BIN="${BIN}" build

clean:
	rm -Rf bin/*

.PHONY: all run build clean
