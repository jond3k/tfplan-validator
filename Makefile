IS_DEV=$(shell git diff --exit-code >/dev/null 2>/dev/null 1>&2 || echo yes)

install: mod test
	go install ./cmd/tfplan-validator

test: mod
	gotestsum --format=short-verbose $(TEST) $(TESTARGS)

coverage: mod
	gotestsum --format=short-verbose -- . ./internal/app/tfplan-validator -coverprofile=coverage.txt -covermode=atomic

coverage-html: coverage
	go tool cover -html=coverage.txt

mod:
	go install gotest.tools/gotestsum@latest
	go mod download && go mod verify && go mod tidy

release-check:
	$(if $(call equals,0,$(shell git diff-index --quiet HEAD; echo $$?)),, \
				$(error Cannot make a release if there are uncommitted changes $?) \
		)

release: mod test release-check
	echo git tag ${RELEASE}
	echo git push origin ${RELEASE}

.PHONY: build install test coverage coverage-html mod release release-check
