.PHONY: build build-run

build:
	@go build -o ./bin/gomerge ./cmd/gomerge

run:
	@./bin/gomerge list -t ${GITHUB_TOKEN} -r cian911/pr-test

build-run: build
	@./bin/gomerge list -t ${GITHUB_TOKEN} -r cian911/pr-test
