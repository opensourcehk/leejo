build: check

check: napping
	go fmt
	go build -o ${BIN}/integration_test

napping: github.com/jmcvetta/napping

github.com/jmcvetta/napping:
	go get github.com/jmcvetta/napping

$PHONY: build napping check
