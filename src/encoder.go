package main

import (
	"encoding/json"
	"encoding/xml"
)

// An Encoder implements an encoding format of values to be sent
// as response to requests on the API endpoints.
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

type JsonEncoder struct {
}

func (enc JsonEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

type XmlEncoder struct {
}

func (enc XmlEncoder) Encode(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}
