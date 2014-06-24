
GOPATH=$(shell pwd)
BIN=${GOPATH}/bin

all: build

run: build

build: bin/leejo_server

bin/leejo_server:
	cd src; make GOPATH="${GOPATH}" BIN="${BIN}" build

clean:
	rm -Rf bin/*

.PHONY: all run build clean
