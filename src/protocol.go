package main

import (
	"encoding/json"
	"io"
	"leejo/rest"
)

type Protocol struct {
}

func (p *Protocol) NewEncoder(w io.Writer) rest.Encoder {
	return json.NewEncoder(w)
}

func (p *Protocol) NewDecoder(r io.Reader) rest.Decoder {
	return json.NewDecoder(r)
}
