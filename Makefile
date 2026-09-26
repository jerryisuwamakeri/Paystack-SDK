.PHONY: build test test-race vet fmt fmt-check lint check

build:
	go build ./...

test:
	go test ./...

test-race:
	go test ./... -race -cover

vet:
	go vet ./...

fmt:
	gofmt -w .

fmt-check:
	@diff <(gofmt -l .) <(echo -n)

check: build vet fmt-check test-race
