package kec

import (
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type ClientKecConfig struct {
	ksyun.ClientConfig `mapstructure:",squash"`
	client             *ClientKecWrapper
}

func (c *ClientKecConfig) kecClient(stateBag *multistep.BasicStateBag) *ClientKecWrapper {
	if c.client != nil {
		return c.client
	}
	c.client = &ClientKecWrapper{
		ClientWrapper: c.Client(stateBag),
		KecClient: kec.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
			&ksc.Config{Region: &c.KsyunRegion},
			&utils.UrlInfo{
				UseSSL: true,
			}),
	}
	return c.client
}
