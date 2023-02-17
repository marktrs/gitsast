.PHONY: clean
clean:
	\rm -rf bin/*

build: clean
	go build -o ./bin/gitsast main.go

run:
	go run main.go

start:
	./bin/gitsast

build-all: clean
	GOOS=linux GOARCH=amd64 go build -o ./bin/gitsast-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o ./bin/gitsast-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o ./bin/gitsast-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o ./bin/gitsast-darwin-arm64 main.go

deps-cleancache:
	go clean -modcache
