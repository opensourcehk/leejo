build: martini martini-render pq
	go fmt
	go build -o ${BIN}/server

martini: github.com/go-martini/martini

pq: github.com/lib/pq

martini-render: github.com/martini-contrib/render

github.com/go-martini/martini:
	go get github.com/go-martini/martini

github.com/martini-contrib/render:
	go get github.com/martini-contrib/render

github.com/lib/pq:
	go get github.com/lib/pq

.PHONY: build margin pq
