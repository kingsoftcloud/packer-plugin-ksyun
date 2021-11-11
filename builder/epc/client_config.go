package epc

import (
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/epc"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type ClientConfig struct {
	*ksyun.AccessConfig
	client *ClientWrapper
}

func (c *ClientConfig) Client() *ClientWrapper {
	if c.client != nil {
		return c.client
	}
	c.client = &ClientWrapper{}
	c.client.EpcClient = epc.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
		&ksc.Config{Region: &c.KsyunRegion}, &utils.UrlInfo{
			UseSSL: true,
		})
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
	return c.client
}
