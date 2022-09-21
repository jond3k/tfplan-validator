build: mod
	go build -o ./bin/tfplan_validator cmd/tfplan_validator/main.go

test: mod
	gotestsum --format=short-verbose $(TEST) $(TESTARGS)

coverage: mod
	gotestsum --format=short-verbose -- . ./cmd -coverprofile=cover.out
	go tool cover -html=cover.out

clear:
	rm bin/tfplan_validator
	rm test-results/*.json

mod:
	go install gotest.tools/gotestsum@latest
	go mod download && go mod verify && go mod tidy

.PHONY: build test coverage mod 
