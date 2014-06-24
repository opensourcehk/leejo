package main

type Resp struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}
