package main

import (
	"log"
	"net/http"
)

func TempServe(addr string, path string, ch chan interface{}, h http.HandlerFunc) {
	http.Handle(path, h)
	log.Fatal(http.ListenAndServe(addr, nil))
}
