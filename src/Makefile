build: martini pq
	go fmt
	go build -o ${BIN}/server

martini: github.com/go-martini/martini

pq: github.com/lib/pq

github.com/go-martini/martini:
	go get github.com/go-martini/martini

github.com/lib/pq:
	go get github.com/lib/pq
