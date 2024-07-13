all: run-app linter-golangci swagger clean

run: run-app

build:
	go build ./cmd/api/main.go
.PHONY: build

run-app:
	cd cmd/api && go mod tidy && go mod download && \
	CGO_ENABLED=0 go run main.go
.PHONY: run-server

swagger:
	swag init -g ./cmd/api/main.go -o ./docs
.PHONY: swagger

linter-golangci:
	golangci-lint run
.PHONY: linter-golangci

clean:
	go clean
.PHONY: clean