GOPATH ?= $(shell go env GOPATH)

.PHONY: lint
lint:
	bash scripts/install-lint.sh
	${GOPATH}/bin/golangci-lint run

test:
	go test ./...

build:
	cd examples && go build -o ../bin/example