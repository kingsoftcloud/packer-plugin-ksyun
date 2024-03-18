This plugin can be used with HashiCorp Packer plugin for creating Kingsoft Cloud KEC & BareMetal image.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

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

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/kingsoftcloud/ksyun
```

### Components

#### Builders
- [ksyun-epc](/packer/integrations/kingsoftcloud/ksyun/latest/components/builder/epc) - The ksyun-epc Packer builder is able to create Ksyun Images backed by
  bare metal instance storage as the root device.
- [ksyun-kec](/packer/integrations/kingsoftcloud/ksyun/latest/components/builder/kec) - The ksyun-kec Packer builder is able to create Ksyun Images backed by
  instance storage as the root device.

#### Data sources
- [ksyun-kmi](/packer/integrations/kingsoftcloud/ksyun/latest/components/data-source/kmi) - The Ksyun KMI data source provides information from a KMI that will be fetched based
  on the filter options provided in the configuration.
