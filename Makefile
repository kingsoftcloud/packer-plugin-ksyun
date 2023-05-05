

HASHICORP_PACKER_PLUGIN_SDK_VERSION?=$(shell go list -m github.com/hashicorp/packer-plugin-sdk | cut -d " " -f2)

default: build

test:
	go test -v ./...

test_integration:
	PACKER_ACC=1 go test -count 1 -v ./...  -timeout=120m

install-packer-sdc: ## Install packer sofware development command
	go install github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc@${HASHICORP_PACKER_PLUGIN_SDK_VERSION}


lint:
	go vet .
	golint .

build:
	@chmod +x scripts/build.sh
	@bash ./scripts/build.sh $(version)

install: build

generate: install-packer-sdc
	go generate ./...
	packer-sdc renderdocs -src ./docs -dst ./.docs -partials ./docs-partials


.PHONY: default test test_integration lint build install
