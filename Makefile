.PHONY: build
build: dist/rekki

dist/rekki: pkg/**/*.go cmd/main.go
	go build -o dist/rekki cmd/main.go

.PHONY: bootstrap
bootstrap:
	go get -u github.com/gorilla/mux