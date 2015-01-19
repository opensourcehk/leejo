package rest

import (
	"github.com/gourd/session"
	"io"
)

// generic protocol interface
type Protocol interface{
	// wrap the response with any protocol structure
	WrapResponse(s session.Session, r interface{}) interface{}

	// wrap the response with any error structure
	WrapError(s session.Session, e error) interface{}

	// switch the encoder with session
	NewEncoder(s session.Session, w io.Writer) Encoder

	// switch the decoder with session
	NewDecoder(s session.Session, r io.Reader) Decoder
}

// generic encoder interface
type Encoder interface{
	Encode(interface{}) error
}

// generic decoder interface
type Decoder interface{
	Decode(interface{}) error
}
