PACKAGES=$(shell go list ./... | grep -v '/vendor/')
BUILD_NUMBER ?= 0-local
TESTNET_NAME ?= localnet
FROM_KEY ?= cosmosaccaddr1kje2wjc66mc3u283dy80czej8m9su8ca5a8drz
AMOUNT ?= 1node101Token
NODE ?= 18.144.38.59:26657
BUILD_FLAGS = -tags netgo -ldflags "-X github.com/greg-szabo/f11/defaults.Release=${BUILD_NUMBER} -X github.com/greg-szabo/f11/defaults.TestnetName=${TESTNET_NAME} -X github.com/greg-szabo/f11/defaults.FromKey=${FROM_KEY} -X github.com/greg-szabo/f11/defaults.Amount=${AMOUNT} -X github.com/greg-szabo/f11/defaults.Node=${NODE}"


########################################
### Build

build-dev:
	go build $(BUILD_FLAGS) -o build/f11 .

build:
	GOOS=linux GOARCH=amd64 $(MAKE) build-dev

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
### Testing

test: test_unit

test_cli:
	@go test -count 1 -p 1 `go list github.com/greg-szabo/f11`

test_unit:
	@go test $(PACKAGES)


########################################
### Localnet (Requirements: pip3 install aws-sam-cli)

localnet-start:
	sam local start-api


########################################
### Release management (set up requirements manually)

package:
	zip "build/f11_${TESTNET_NAME}.zip" build/f11 template.yml
	aws s3 cp "build/f11_${TESTNET_NAME}.zip" "s3://tendermint-lambda/f11_${TESTNET_NAME}.zip"

deploy:
	sam deploy --template-file template.yml --stack-name "f11_${TESTNET_NAME}"" --capabilities CAPABILITY_IAM --region us-east-1

.PHONY: build build-dev check_tools update_tools get_tools get_vendor_deps test test_cli test_unit package deploy

