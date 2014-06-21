
GOPATH=$(shell pwd)
BIN=${GOPATH}/bin

all: build

build:
	cd src; make GOPATH="${GOPATH}" BIN="${BIN}" $@

