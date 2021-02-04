include .bingo/Variables.mk
FILES_TO_FMT      ?= $(shell find . -path ./vendor -prune -o -name '*.go' -print)

GO111MODULE       ?= on
export GO111MODULE

GOBIN             ?= $(firstword $(subst :, ,${GOPATH}))/bin

DOCKER_IMAGE_REPO ?= quay.io/prmsrswt/foundry
DOCKER_IMAGE_TAG  ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))-$(shell date +%Y-%m-%d)-$(shell git rev-parse --short HEAD)

help: ## Displays help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: all
all: format build

.PHONY: build
build: ## Build foundry.
	@echo ">> building foundry"
	@go build ./cmd/foundry

.PHONY: install
install: ## Build and install foundry.
	@echo ">> installing foundry"
	@go install github.com/prmsrswt/foundry/cmd/foundry

.PHONY: deps
deps: ## Ensures fresh go.mod and go.sum.
	@go mod tidy
	@go mod verify

.PHONY: format
format: ## Formats Go code including imports.
format: $(GOIMPORTS)
	@echo ">> formatting code"
	@$(GOIMPORTS) -w $(FILES_TO_FMT)

.PHONY: test
test: ## Runs all Go unit tests.
	@echo ">> running unit tests"
	@go test -v -timeout=30m $(shell go list ./... | grep -v /vendor/)

.PHONY: coverage
coverage: ## Create test coverage profile
	@echo ">> creating test coverage profile"
	@go test -v -covermode=count -coverprofile=coverage.out $(shell go list ./... | grep -v /vendor/)

.PHONY: lint
lint: ## Runs various static analysis against our code.
lint: $(GOLANGCI_LINT) format deps
	@echo ">> examining all of the Go files"
	@go vet -stdmethods=false ./...
	@echo ">> linting all of the Go files GOGC=${GOGC}"
	@$(GOLANGCI_LINT) run

.PHONY: docker
docker: ## Build Docker image
	@echo ">> building docker image 'foundry'"
	@docker build -t foundry .

.PHONY: docker-push
docker-push: ## Pushes 'foundry' docker image to "$(DOCKER_IMAGE_REPO):$(DOCKER_IMAGE_TAG)".
	@echo ">> pushing image"
	@docker tag "foundry" "$(DOCKER_IMAGE_REPO):$(DOCKER_IMAGE_TAG)"
	@docker push "$(DOCKER_IMAGE_REPO):$(DOCKER_IMAGE_TAG)"
