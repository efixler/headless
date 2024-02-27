MODULE_NAME := $(shell go list -m)
BUILD_DIR := build

.DEFAULT_GOAL := build

.PHONY: fmt vet build clean

build: vet ## build the binaries, to the build/ folder (default target)
	@echo "Building $(MODULE_NAME)..."
	@go build -o $(BUILD_DIR)/ -tags "$(TAGS)" ./...
	
clean: ## clean the build directory
	@echo "Cleaning $(MODULE_NAME)..."
	@rm -rf $(BUILD_DIR)/*

docker-pull-headless-chrome:
	@echo "Pulling headless chrome..."
	@docker pull yukinying/chrome-headless-browser-stable

run-headless-chrome: docker-pull-headless-chrome ## run headless chrome
	@echo "Running headless chrome on port 9222..."
	@docker run --init -it --rm --name headless-chrome -p=127.0.0.1:9222:9222 --shm-size=1024m --cap-add=SYS_ADMIN yukinying/chrome-headless-browser-stable

fmt: 
	@echo "Running go fmt..."
	@go fmt ./...

test: ## run the tests
	@echo "Running tests..."
	@go test -coverprofile=coverage.out ./... 

vet: fmt ## fmt, vet, and staticcheck
	@echo "Running go vet and staticcheck..."
	@go vet ./...
	@staticcheck ./...

cognitive: ## run the cognitive complexity checker
	@echo "Running gocognit..."
	@gocognit  -ignore "_test|testdata" -top 5 .

help: ## show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)