package rest

import (
	"encoding/json"
	"github.com/gourd/service"
	"leejo/data"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, err error) {

	if err == nil {
		return
	}

	resp := data.Resp{
		Status: "fail",
		Code:   500,
	}

	if se, ok := err.(service.EntityError); ok {
		resp.Code = se.Code()
		resp.Message = se.Error()
	} else {
		resp.Message = "Internal Server Error"
	}

	log.Printf("Internal Server Error: %s", err.Error())

	w.WriteHeader(resp.Code)
	json.NewEncoder(w).Encode(resp)

}
