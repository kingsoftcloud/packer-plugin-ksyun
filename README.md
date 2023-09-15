# Packer Builder for Kingsoft Cloud KEC And Bare Metal

This is a [HashiCorp Packer](https://www.packer.io/) plugin for creating [Kingsoft Cloud KEC & BareMetal](https://www.ksyun.com/nv/product/KEC.html) image.

## Requirements
* [Go 1.19+](https://golang.org/doc/install)
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
$ make install
```

### Install from HCL:
```hcl
packer {
  required_plugins {
    ksyun = {
      version = ">= 0.0.10"
      source  = "github.com/kingsoftcloud/ksyun"
    }
  }
}
```


### Install from release:

* Download binaries from the [releases page](https://github.com/kingsoftcloud/packer-plugin-ksyun/releases).
* [Install](https://www.packer.io/docs/extending/plugins.html#installing-plugins) the plugin, or simply put it into the same directory with JSON templates.
* Move the downloaded binary to `~/.packer.d/plugins/`

## Usage for Kec
Here is a sample template, which you can also find in the `example/` directory
### JSON Template
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
### HCL Template

```hcl
variable "access_key" {
  type    = string
  default = env("KSYUN_ACCESS_KEY")
}

variable "secret_key" {
  type    = string
  default = env("KSYUN_SECRET_KEY")
}

source "ksyun-kec" "example" {
  access_key                  = var.access_key
  secret_key                  = var.secret_key
  region                      = "cn-shanghai-2"
  image_name                  = "packer_test"
  source_image_id             = "IMG-dd1f8324-1f27-46e0-ad6b-b41d8c8ff025"
  instance_type               = "N3.1B"
  ssh_username                = "root"
  associate_public_ip_address = true
}
build {
  provisioner "shell" {
    inline = [
      "sleep 30",
      "yum install mysql -y"
    ]
  }
}

```
## Usage for Bare Metal
Here is a sample template, which you can also find in the `example/` directory

### JSON Template

```json
{
  "variables": {
    "access_key": "{{ env `KSYUN_ACCESS_KEY` }}",
    "secret_key": "{{ env `KSYUN_SECRET_KEY` }}"
  },
  "builders": [{
    "type":"ksyun-epc",
    "access_key":"{{user `access_key`}}",
    "secret_key":"{{user `secret_key`}}",
    "region":"cn-beijing-6",
    "source_image_id":"eb8c0428-476e-49af-8ccb-9fad2455a54c",
    "host_type":"EC-I-III-II",
    "availability_zone":"cn-beijing-6c",
    "raid": "Raid1",
    "ssh_username":"root",
    "ssh_clear_authorized_keys": true,
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

### HCL Template

```hcl
variable "access_key" {
  type    = string
  default = env("KSYUN_ACCESS_KEY")
}

variable "secret_key" {
  type    = string
  default = env("KSYUN_SECRET_KEY")
}
source "ksyun-kec" "example" {
  access_key                  = var.access_key
  secret_key                  = var.secret_key
  region                      = "cn-beijing-6"
  source_image_id             = "eb8c0428-476e-49af-8ccb-9fad2455a54c"
  host_type                   = "EC-I-III-II"
  availability_zone           = "cn-beijing-6c"
  raid                        = "Raid1"
  ssh_username                = "root"
  ssh_clear_authorized_keys   = true
  associate_public_ip_address = true
}
build {
  provisioner "shell" {
    inline = [
      "sleep 30",
      "yum install mysql -y"
    ]
  }
}

```

Enter the API user credentials in your terminal with the following commands. Replace the <AK> and <SK> with your user details.
```sh
export KSYUN_ACCESS_KEY=<AK>
export KSYUN_SECRET_KEY=<SK>
```
Then run Packer using the example template with the command underneath.

### HTL Template

```
# install packer plugin
packer init .
# use for kec
packer build -only="ksyun.*" .
# use for bare metal
packer build -only="ksyun_epc.*" .
```

### JSON Template
```
# use for kec
packer build example/ksyun.json
# use for bare metal
packer build example/ksyun_epc.json
```


