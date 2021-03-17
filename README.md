# Packer Builder for Kingsoft Cloud KEC

This is a [HashiCorp Packer](https://www.packer.io/) plugin for creating [Kingsoft Cloud KEC](https://www.ksyun.com/nv/product/KEC.html) image.

## Requirements
* [Go 1.14+](https://golang.org/doc/install)
* [Packer](https://www.packer.io/intro/getting-started/install.html)

## Build & Installation

### Install from source:

Clone repository to `$GOPATH/src/github.com/kingsoftcloud/packer-plugin-ksyun`

```sh
$ mkdir -p $GOPATH/src/github.com/kingsoftcloud; 
$ cd $GOPATH/src/github.com/kingsoftcloud
$ git clone git@github.com:kingsoftcloud/packer-plugin-ksyun.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/kingsoftcloud/packer-plugin-ksyun
$ make build
```

Link the build to Packer

```sh
$ln -s $GOPATH/bin/packer-plugin-ksyun ~/.packer.d/plugins/packer-plugin-ksyun
```

### Install from release:

* Download binaries from the [releases page](https://github.com/kingsoftcloud/packer-plugin-ksyun/releases).
* [Install](https://www.packer.io/docs/extending/plugins.html#installing-plugins) the plugin, or simply put it into the same directory with JSON templates.
* Move the downloaded binary to `~/.packer.d/plugins/`

