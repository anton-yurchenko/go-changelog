I := "⚪"
E := "🔴"
BINARY := $(notdir $(CURDIR))
GO_BIN_DIR := $(GOPATH)/bin

test: lint
	@echo "$(I) unit testing..."
	@go test $$(go list ./... | grep -v vendor | grep -v mocks) -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: lint
lint: $(GO_LINTER)
	@echo "$(I) vendoring..."
	@go mod vendor || (echo "$(E) 'go mod vendor' error"; exit 1)
	@go mod tidy || (echo "$(E) 'go mod tidy' error"; exit 1)
	@echo "$(I) linting..."
	@golangci-lint run ./... || (echo "$(E) linter error"; exit 1)

.PHONY: init
init:
	@echo "$(I) initializing project..."
	@rm -f go.mod
	@rm -f go.sum
	@rm -rf ./vendor
	@go mod init $$(pwd | awk -F'/' '{print "github.com/"$$(NF-1)"/"$$NF}')

GO_LINTER := $(GO_BIN_DIR)/golangci-lint
$(GO_LINTER):
	@echo "$(I) installing linter..."
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

.PHONY: codecov
codecov: test
	@go tool cover -html=coverage.txt || (echo "$(E) 'go tool cover' error"; exit 1)
