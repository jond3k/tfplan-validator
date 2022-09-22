build: mod
	mkdir -p ./bin
	go build -o ./bin/tfplan-validator ./cmd/tfplan-validator/

install: mod
	go install ./cmd/tfplan-validator/

test: mod
	mkdir -p ./test-results
	gotestsum --format=short-verbose $(TEST) $(TESTARGS)

coverage: mod
	mkdir -p ./test-results
	gotestsum --format=short-verbose -- . ./cmd -coverprofile=coverage.txt -covermode=atomic

coverage-html: coverage
	go tool cover -html=coverage.txt

mod:
	go install gotest.tools/gotestsum@latest
	go mod download && go mod verify && go mod tidy

.PHONY: build install test coverage coverage-html mod
