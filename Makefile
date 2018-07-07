PACKAGES=$(shell go list ./... | grep -v '/vendor/')
BUILD_NUMBER ?= 0-local
TESTNET_NAME ?= localnet
BUILD_FLAGS = -tags netgo -ldflags "-X github.com/greg-szabo/f11/version.Release=${BUILD_NUMBER} -X github.com/greg-szabo/f11/testnet.TestnetName=${TESTNET_NAME}"


########################################
### Build

build:
	go build $(BUILD_FLAGS) -o build/f11 .

build-linux:
	GOOS=linux GOARCH=amd64 $(MAKE) build

########################################
### Tools & dependencies

DEP = github.com/golang/dep/cmd/dep
DEP_CHECK := $(shell command -v dep 2> /dev/null)

check_tools:
	cd tools && $(MAKE) check_tools

update_tools:
	cd tools && $(MAKE) update_tools

get_tools:
	cd tools && $(MAKE) get_tools

get_vendor_deps:
	@rm -rf vendor/
	@echo "--> Running dep ensure"
	@dep ensure -v


########################################
### Localnet

localnet-start: build-linux
	sam local start-api


########################################
### Testing

test: test_unit

test_cli:
	@go test -count 1 -p 1 `go list github.com/greg-szabo/f11`

test_unit:
	@go test $(PACKAGES)

.PHONY: build check_tools update_tools get_tools get_vendor_deps test test_cli test_unit
