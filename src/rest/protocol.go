package rest

import (
	"io"
)

// generic protocol interface
type Protocol interface{
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

// generic encoder interface
type Encoder interface{
	Encode(interface{}) error
}

// generic decoder interface
type Decoder interface{
	Decode(interface{}) error
}
