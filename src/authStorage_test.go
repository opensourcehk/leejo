package main

import (
	"github.com/RangelReale/osin"
	"testing"
)

// just test if the passed in parameteer
// implements osin.Stroage
func useOsinStorage(s osin.Storage) {
}

func Test_AuthStorage(t *testing.T) {
	a := AuthStorage{}
	useOsinStorage(&a)
}
