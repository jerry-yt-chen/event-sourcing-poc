PUB_NAME=event-publisher
SUB_NAME=event-subscriber
VERSION=latest

.PHONY: run
run: build
	docker-compose -f deployments/docker-compose.yaml up

.PHONY: build
build:
	docker build -t $(PUB_NAME):$(VERSION) -f ./build/$(PUB_NAME)/Dockerfile .
	docker build -t $(SUB_NAME):$(VERSION) -f ./build/$(SUB_NAME)/Dockerfile .

.PHONY: down
down:
	docker-compose -f deployments/docker-compose.yaml down

.PHONY: di
di:
	go install github.com/google/wire/cmd/wire@v0.5.0
	wire gen ./internal/injector
