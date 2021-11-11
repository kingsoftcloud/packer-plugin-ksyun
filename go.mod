module github.com/kingsoftcloud/packer-plugin-ksyun

go 1.15

require (
	github.com/KscSDK/ksc-sdk-go v0.1.41
	github.com/hashicorp/hcl/v2 v2.8.0
	github.com/hashicorp/packer-plugin-sdk v0.0.14
	github.com/zclconf/go-cty v1.7.0
)

replace github.com/KscSDK/ksc-sdk-go => ../../KscSDK/ksc-sdk-go