.PHONY: build build-run
.PHONY: test-all test-cli test-client test-utils

PKG := ./

VERSION := valpha-test.1
BUILD := $$(git log -1 --pretty=%h)
BUILD_DATE := $$(date -u +"%Y%m%d.%H%M%S")

build:
	@go build -o ./bin/gomerge ./cmd/gomerge
	@go build \
		-ldflags "-X main.Version=${VERSION} \
							-X main.Build=${BUILD} \
							-X main.BuildDate=${BUILD_DATE}" \
		-o ./bin/gomerge ./cmd/gomerge

run:
	@./bin/gomerge list -t ${GITHUB_TOKEN} -r cian911/pr-test

build-run: build
	@./bin/gomerge list -t ${GITHUB_TOKEN} -r cian911/pr-test

test-all: test-cli test-client test-utils

test-cli:
	@go test -v ${PKG}/pkg/cli/list

test-client:
	@go test -v ${PKG}/pkg/gitclient

test-utils:
	@go test -v ${PKG}/pkg/utils
