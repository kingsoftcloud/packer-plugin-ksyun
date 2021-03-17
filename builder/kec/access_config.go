//go:generate struct-markdown

package kec

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"os"
)

type KsyunAccessConfig struct {
	// Ksyun access key must be provided unless `profile` is set, but it can
	// also be sourced from the `KSYUN_ACCESS_KEY` environment variable.
	KsyunAccessKey string `mapstructure:"access_key" required:"true"`
	// Ksyun secret key must be provided unless `profile` is set, but it can
	// also be sourced from the `KSYUN_SECRET_KEY` environment variable.
	KsyunSecretKey string `mapstructure:"secret_key" required:"true"`
	// Ksyun region must be provided unless `profile` is set, but it can
	// also be sourced from the `KSYUN_REGION` environment variable.
	KsyunRegion string `mapstructure:"region" required:"true"`

	client *ClientWrapper
}

func (c *KsyunAccessConfig) Client() *ClientWrapper {
	if c.client != nil {
		return c.client
	}
	c.client = &ClientWrapper{}
	c.client.KecClient = kec.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
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

func (c *KsyunAccessConfig) Config() error {
	if c.KsyunAccessKey == "" {
		c.KsyunAccessKey = os.Getenv("KSYUN_ACCESS_KEY")
	}
	if c.KsyunSecretKey == "" {
		c.KsyunSecretKey = os.Getenv("KSYUN_SECRET_KEY")
	}
	if c.KsyunAccessKey == "" || c.KsyunSecretKey == "" {
		return fmt.Errorf("KSYUN_ACCESS_KEY and KSYUN_SECRET_KEY must be set in template file or environment variables")
	}
	return nil

}

func (c *KsyunAccessConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	if err := c.Config(); err != nil {
		errs = append(errs, err)
	}

	if c.KsyunRegion == "" {
		c.KsyunRegion = os.Getenv("KSYUN_REGION")
	}

	if c.KsyunRegion == "" {
		errs = append(errs, fmt.Errorf("region option or KSYUN_REGION must be provided in template file or environment variables"))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}
