MAKEFLAGS += --no-print-directory
GOBIN ?= $(shell go env GOPATH)/bin

deps:
	go mod download -x

testdeps: deps
	go install honnef.co/go/tools/cmd/staticcheck@2022.1.3

tidy:
	go mod verify
	go mod tidy

vet: testdeps
	go vet ./...

staticcheck: testdeps
	$(GOBIN)/staticcheck ./...

lint: vet staticcheck

test:
	go test -v -covermode=atomic -coverprofile=coverage.out ./...

check: test lint

clean:
	go clean ./...
