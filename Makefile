.PHONY: build test vet fmt-check clean

build:
	go build -o bin/authopus ./cmd/authopus
	go build -o bin/authopusctl ./cmd/authopusctl

test:
	go test ./...

vet:
	go vet ./...

fmt-check:
	@unformatted=$$(gofmt -l .) || exit $$?; \
	if [ -n "$$unformatted" ]; then \
		printf 'Go files need gofmt:\n%s\n' "$$unformatted"; \
		exit 1; \
	fi

clean:
	rm -rf bin
