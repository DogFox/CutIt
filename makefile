
BINARY=previewer
DOCKER_IMAGE=previewer:latest

build:
	go build -o $(BINARY) ./cmd/main.go

docker-build:
	docker build -t $(DOCKER_IMAGE) .

run:
	docker-compose up --build -d

test:
	go test ./...
	go test -v ./tests/integration_test.go

clean:
	rm -f $(BINARY)
	docker-compose down

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.63.4

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run test clean lint
