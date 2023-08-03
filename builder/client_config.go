package ksyun

import (
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/tagv2"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type ClientConfig struct {
	AccessConfig `mapstructure:",squash"`
	client       *ClientWrapper
}

func (c *ClientConfig) Client(stateBag *multistep.BasicStateBag) *ClientWrapper {
	if c.client != nil {
		stateBag.Put("ksyun_client", c.client)
		return c.client
	}
	c.client = &ClientWrapper{}
	c.client.SksClient = sks.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
		&ksc.Config{Region: &c.KsyunRegion}, &utils.UrlInfo{
			UseSSL: true,
		})
	c.client.EipClient = eip.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
		&ksc.Config{Region: &c.KsyunRegion}, &utils.UrlInfo{
			UseSSL: true,
		})
	c.client.VpcClient = vpc.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
		&ksc.Config{Region: &c.KsyunRegion}, &utils.UrlInfo{
			UseSSL: true,
		})

	c.client.TagsClient = tagv2.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
		&ksc.Config{Region: &c.KsyunRegion}, &utils.UrlInfo{
			UseSSL: true,
		})
	stateBag.Put("ksyun_client", c.client)
	return c.client
}
