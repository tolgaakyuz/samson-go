GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_BUILD_RACE=$(GO_CMD) build -race
GO_TEST=$(GO_CMD) test
GO_TEST_VERBOSE=$(GO_CMD) test -v
GO_INSTALL=$(GO_CMD) install -v
GO_CLEAN=$(GO_CMD) clean
GO_DEPS=$(GO_CMD) get -d -v
GO_DEPS_UPDATE=$(GO_CMD) get -d -v -u
GO_VET=$(GO_CMD) vet
GO_FMT=$(GO_CMD) fmt

.PHONY: all test lint dep build

all: build

build-test:
	docker-compose pull samson
	docker-compose build --pull samson-go-tests

test: build-test
	DOCKER_REGISTRY=none docker-compose run \
					-e TESTDATA=/go/src/github.com/tolgaakyuz/samson-go/testdata \
					samson-go-tests \
					bash -c './tools/test.sh'
test-i: build-test
	DOCKER_REGISTRY=none docker-compose run \
		-e TESTDATA=/go/src/github.com/tolgaakyuz/samson-go/testdata \
		samson-go-tests \
		bash

coverage:
	go tool cover -html=coverage.txt

lint:
	./tools/lint.sh
