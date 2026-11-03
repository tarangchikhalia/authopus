.PHONY: build test vet fmt-check clean

build:
	go build -o bin/authopus ./cmd/authopus
	go build -o bin/authopusctl ./cmd/authopusctl

test:
	go test ./...

vet:
	go vet ./...

fmt-check:
	@test -z "$$(gofmt -l $$(find cmd internal -name '*.go'))" || (echo "Go files need gofmt"; exit 1)

clean:
	rm -rf bin
