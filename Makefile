
GOPATH=$(shell pwd)
BIN=${GOPATH}/bin

all: build

run: build
	${BIN}/server

build:
	cd src; make GOPATH="${GOPATH}" BIN="${BIN}" $@

.PHONY: all run build
