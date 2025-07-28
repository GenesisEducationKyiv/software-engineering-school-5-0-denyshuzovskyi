LINT_CONFIG=./../.golangci.yaml

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run --config $(LINT_CONFIG)

.PHONY: fmt
fmt:
	golangci-lint fmt --config $(LINT_CONFIG)

.PHONY: test
test:
	go test ./...

.PHONY: testit
testit:
	go test -race --tags=integration ./...
