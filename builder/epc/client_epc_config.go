package epc

import (
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/epc"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type ClientEpcConfig struct {
	ksyun.ClientConfig `mapstructure:",squash"`
	client             *ClientEpcWrapper
}

func (c *ClientEpcConfig) EpcClient(stateBag *multistep.BasicStateBag) *ClientEpcWrapper {
	if c.client != nil {
		return c.client
	}
	c.client = &ClientEpcWrapper{
		ClientWrapper: c.Client(stateBag),
		EpcClient: epc.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
			&ksc.Config{Region: &c.KsyunRegion},
			&utils.UrlInfo{
				UseSSL: true,
			}),
	}
	return c.client
}
