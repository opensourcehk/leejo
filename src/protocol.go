package main

import (
	"encoding/json"
	"github.com/gourd/service"
	"github.com/gourd/session"
	"io"
	"leejo/data"
	"leejo/rest"
	"log"
)

type Protocol struct {
}

func (p *Protocol) Response(s session.Session, r interface{}) interface{} {
	return &data.Resp{
		Status: "OK",
		Result: r,
	}
}

func (p *Protocol) Error(s session.Session, err error) interface{} {
	if err == nil {
		return nil
	}

	resp := data.Resp{
		Status: "fail",
		Code:   500,
	}

	if se, ok := err.(service.EntityError); ok {
		log.Printf("error is service.EntityError")
		resp.Code = se.Code()
		resp.Message = se.Error()
	} else {
		log.Printf("error is not service.EntityError")
		resp.Message = "Internal Server Error"
	}

	return &resp
}

func (p *Protocol) NewEncoder(s session.Session, w io.Writer) rest.Encoder {
	return json.NewEncoder(w)
}

func (p *Protocol) NewDecoder(s session.Session, r io.Reader) rest.Decoder {
	return json.NewDecoder(r)
}
