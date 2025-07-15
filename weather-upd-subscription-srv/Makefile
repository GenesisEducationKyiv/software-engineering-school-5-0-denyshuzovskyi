LINT_CONFIG=./../.golangci.yaml

.PHONY: tidy lint fmt test testit build-img

tidy:
	go mod tidy

lint:
	golangci-lint run --config $(LINT_CONFIG)

fmt:
	golangci-lint fmt --config $(LINT_CONFIG)

test:
	go test ./...

testit:
	go test --tags=integration ./...
