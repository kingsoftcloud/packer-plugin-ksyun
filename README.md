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

## Usage
Here is a sample template, which you can also find in the `example/` directory
```json
{
  "variables": {
    "access_key": "{{ env `KSYUN_ACCESS_KEY` }}",
    "secret_key": "{{ env `KSYUN_SECRET_KEY` }}"
  },
  "builders": [{
    "type":"ksyun-kec",
    "access_key":"{{user `access_key`}}",
    "secret_key":"{{user `secret_key`}}",
    "region":"cn-shanghai-2",
    "image_name":"packer_test",
    "source_image_id":"IMG-dd1f8324-1f27-46e0-ad6b-b41d8c8ff025",
    "instance_type":"N3.1B",
    "ssh_username":"root",
    "associate_public_ip_address": true
  }],
  "provisioners": [{
    "type": "shell",
    "inline": [
      "sleep 30",
      "yum install mysql -y"
    ]
  }]
}
```
Enter the API user credentials in your terminal with the following commands. Replace the <AK> and <SK> with your user details.
```sh
export KSYUN_ACCESS_KEY=<AK>
export KSYUN_SECRET_KEY=<SK>
```
Then run Packer using the example template with the command underneath.
```
packer build example/ksyun.json
```


