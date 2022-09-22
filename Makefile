build: mod
	go build -o ./bin/tfplan_validator cmd/tfplan_validator/main.go

test: mod
	mkdir -p ./test-results
	gotestsum --format=short-verbose $(TEST) $(TESTARGS)

coverage: mod
	mkdir -p ./test-results
	gotestsum --format=short-verbose -- . ./cmd -coverprofile=coverage.txt -covermode=atomic

coverage-html: coverage
	go tool cover -html=coverage.txt

clear:
	rm bin/tfplan_validator
	rm test-results/*.json

mod:
	go install gotest.tools/gotestsum@latest
	go mod download && go mod verify && go mod tidy

.PHONY: build test coverage mod 
