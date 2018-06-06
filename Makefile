all: build

build-test:
	docker-compose build samson-go-tests

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

.PHONY: all test lint dep build
