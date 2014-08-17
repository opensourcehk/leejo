package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	Db dbconfig `json:"db"`
}

type dbconfig struct {
	Host string `json:"host"`
	Name string `json:"name"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func getConfig() (conf config) {

	// parse flags
	confFn := flag.String("config", "config.json", "Path to config file")
	flag.Parse()

	// read the config file to conf
	confFile, err := ioutil.ReadFile(*confFn)
	if err != nil {
		log.Printf("Failed opening config file \"%s\": %v", *confFn, err)
		os.Exit(1)
	}

	// read the config
	err = json.Unmarshal(confFile, &conf)
	if err != nil {
		log.Printf("Failed parsing config file \"%s\": %v", *confFn, err)
		os.Exit(1)
	}

	return

}
