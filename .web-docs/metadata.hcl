# For full specification on the configuration of this file visit:
# https://github.com/hashicorp/integration-template#metadata-configuration
integration {
  name = "Kingsoft Cloud"
  description = "The Kingsoft Cloud plugin can be used with HashiCorp Packer to create custom images on Ksyun."
  identifier = "packer/kingsoftcloud/ksyun"
  component {
    type = "data-source"
    name = "Ksyun KMI"
    slug = "kmi"
  }
  component {
    type = "builder"
    name = "Ksyun EPC"
    slug = "epc"
  }
  component {
    type = "builder"
    name = "Ksyun KEC"
    slug = "kec"
  }
}
