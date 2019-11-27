default: lint test build

tools: ## Install the tools used to test and build
	@echo "==> Installing build tools"
	GO111MODULE=off go get -u github.com/ahmetb/govvv
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

build: ## Build Chemtrail for development purposes
	@echo "==> Running $@..."
# 	govvv build -o chemtrail ./cmd -version $(shell git describe --tags --abbrev=0 $(git rev-list --tags --max-count=1) |cut -c 2- |awk '{print $1}')+dev -pkg github.com/jrasell/chemtrail/pkg/build
	go build -o chemtrail cmd/main.go

test: ## Run the Chemtrail test suite with coverage
	@echo "==> Running $@..."
	@go test ./... -cover -v -tags -race \
		"$(BUILDTAGS)" $(shell go list ./... | grep -v vendor)

release: ## Trigger the release build script
	@echo "==> Running $@..."
	@goreleaser --rm-dist

.PHONY: lint
lint: ## Run golangci-lint
	@echo "==> Running $@..."
	golangci-lint run cmd/... pkg/...

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
.PHONY: help
help: ## Display this usage information
	@echo "Chemtrail make commands:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
