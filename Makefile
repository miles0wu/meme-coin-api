VERSION := $(shell git describe --tags --always)

all: help

help: ## show help
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'

.PHONY: clean
clean: ## clean artifacts
	@rm -rf ./coverage.txt ./*.out
	@rm -rf ./out ./.bin
	@echo Successfully removed artifacts

.PHONY: version
version: ## show version
	@echo $(VERSION)

.PHONY: dev
dev: ## run dev server
	docker compose up --build

.PHONY: lint
lint: ## run golangci-lint
	@golangci-lint run ./...

.PHONY: test-go
test-go: ## run go test
	@sh $(shell pwd)/script/go.test.sh

.PHONY: gen
gen: gen-source-go gen-swagger ## generate all

.PHONY: gen-source-go
gen-source-go:
	## Starting generate wire and mockgen
	@go generate -tags="wireinject" ./...
	@echo Successfully generated wire and mockgen

.PHONY: gen-swagger
gen-swagger: ## generate swagger
	@swag init -q -o ./api/docs