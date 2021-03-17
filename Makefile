default: build

test:
	go test -v ./...

test_integration:
	PACKER_ACC=1 go test -count 1 -v ./...  -timeout=120m

lint:
	go vet .
	golint .

build:
	go build -v

install: build
	mkdir -p ~/.packer.d/plugins
	install ./packer-plugin-ksyun ~/.packer.d/plugins/

.PHONY: default test test_integration lint build install
