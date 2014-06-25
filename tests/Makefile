build: check

check: napping restit
	go fmt
	go build -o ${BIN}/integration_test

napping: github.com/jmcvetta/napping

restit: github.com/yookoala/restit

github.com/yookoala/restit:
	go get github.com/yookoala/restit

github.com/jmcvetta/napping:
	go get github.com/jmcvetta/napping

$PHONY: build napping restit check
