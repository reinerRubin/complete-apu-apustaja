build-docker:
	docker-compose build
.PHONY: build-docker

run-docker:
	docker-compose up
.PHONY: run-docker

build:
	mkdir -p bin
	CGO_ENABLED=0 go build -o bin/completer cmd/complete-server/complete-server.go
.PHONY: build

run:
	./bin/completer
.PHONY: run
