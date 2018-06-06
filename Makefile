all: dep test lint

test:
	./tools/test.sh

lint:
	./tools/lint.sh

dep:
	dep ensure -vendor-only

.PHONY: all test lint dep
