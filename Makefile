GROUP=kingsoftcloud
SHORT_NAME=ksyun
PLUGIN_NAME=packer-plugin-${SHORT_NAME}


default: build

test:
	go test -v ./...

test_integration:
	PACKER_ACC=1 go test -count 1 -v ./...  -timeout=120m

lint:
	go vet .
	golint .

build:
	@chmod +x scripts/build.sh
	@bash ./scripts/build.sh $(version) $(GROUP) $(SHORT_NAME) $(PLUGIN_NAME)
	#go build -o $(PLUGIN_NAME) -v

install: build
	mkdir -p ~/.packer.d/plugins
	install ./packer-plugin-ksyun ~/.packer.d/plugins/


.PHONY: default test test_integration lint build install
