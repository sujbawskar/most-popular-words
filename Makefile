TEST_FLAGS ?= -race -failfast -count=1
APP_DOCKER_LABEL = cmd-assignment
REVISION = $(shell git rev-parse --short HEAD)
APP_DOCKER_TAG = $(VERSION).$(BUILD_NUMBER)-$(REVISION)

.PHONY: run
run:
	go run .

.PHONY: test
test:
	go test $(TEST_FLAGS) ./...
clean:
	go clean
build:
	go build .

docker.build:
	docker build -t $(APP_DOCKER_LABEL):$(APP_DOCKER_TAG)