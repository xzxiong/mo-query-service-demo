
all: build

clean:
	rm -rf ./bin

.PHONY: fmt
fmt:
	go fmt ./...

build: fmt
	go build -o bin/demo cmd/main.go
