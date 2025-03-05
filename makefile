
BINARY=previewer
DOCKER_IMAGE=previewer:latest

build:
	go build -o $(BINARY) ./cmd/main.go

docker-build:
	docker build -t $(DOCKER_IMAGE) .

run:
	docker-compose up --build

test:
	go test ./...
	go test -v ./tests/integration_test.go

clean:
	rm -f $(BINARY)
	docker-compose down

.PHONY: build run test clean
