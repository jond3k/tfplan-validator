build: mod
	go build -o ./bin/tfplan_validator cmd/tfplan_validator/main.go

test: mod
	gotestsum --format=short-verbose $(TEST) $(TESTARGS)

coverage: mod
	gotestsum --format=short-verbose -- -coverprofile=cover.out
	go tool cover -html=cover.out

mod:
	go install gotest.tools/gotestsum@latest
	go mod download && go mod verify && go mod tidy

.PHONY: build test coverage mod 
