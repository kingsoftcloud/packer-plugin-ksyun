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

type KingcloudAccessConfig struct {
	// Kingcloud access key must be provided unless `profile` is set, but it can
	// also be sourced from the `KINGCLOUD_ACCESS_KEY` environment variable.
	KingcloudAccessKey string `mapstructure:"access_key" required:"true"`
	// Kingcloud secret key must be provided unless `profile` is set, but it can
	// also be sourced from the `KINGCLOUD_SECRET_KEY` environment variable.
	KingcloudSecretKey string `mapstructure:"secret_key" required:"true"`
	// Kingcloud region must be provided unless `profile` is set, but it can
	// also be sourced from the `KINGCLOUD_REGION` environment variable.
	KingcloudRegion string `mapstructure:"region" required:"true"`

	client *ClientWrapper
}

func (c *KingcloudAccessConfig) Client() *ClientWrapper {
	if c.client != nil {
		return c.client
	}
	c.client = &ClientWrapper{}
	c.client.KecClient = kec.SdkNew(ksc.NewClient(c.KingcloudAccessKey, c.KingcloudSecretKey),
		&ksc.Config{Region: &c.KingcloudRegion}, &utils.UrlInfo{
			UseSSL: true,
		})
	c.client.SksClient = sks.SdkNew(ksc.NewClient(c.KingcloudAccessKey, c.KingcloudSecretKey),
		&ksc.Config{Region: &c.KingcloudRegion}, &utils.UrlInfo{
			UseSSL: true,
		})
	c.client.EipClient = eip.SdkNew(ksc.NewClient(c.KingcloudAccessKey, c.KingcloudSecretKey),
		&ksc.Config{Region: &c.KingcloudRegion}, &utils.UrlInfo{
			UseSSL: true,
		})
	c.client.VpcClient = vpc.SdkNew(ksc.NewClient(c.KingcloudAccessKey, c.KingcloudSecretKey),
		&ksc.Config{Region: &c.KingcloudRegion}, &utils.UrlInfo{
			UseSSL: true,
		})
	return c.client
}

func (c *KingcloudAccessConfig) Config() error {
	if c.KingcloudAccessKey == "" {
		c.KingcloudAccessKey = os.Getenv("KINGCLOUD_ACCESS_KEY")
	}
	if c.KingcloudSecretKey == "" {
		c.KingcloudSecretKey = os.Getenv("KINGCLOUD_SECRET_KEY")
	}
	if c.KingcloudAccessKey == "" || c.KingcloudSecretKey == "" {
		return fmt.Errorf("KINGCLOUD_ACCESS_KEY and KINGCLOUD_SECRET_KEY must be set in template file or environment variables")
	}
	return nil

}

func (c *KingcloudAccessConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	if err := c.Config(); err != nil {
		errs = append(errs, err)
	}

	if c.KingcloudRegion == "" {
		c.KingcloudRegion = os.Getenv("KINGCLOUD_REGION")
	}

	if c.KingcloudRegion == "" {
		errs = append(errs, fmt.Errorf("region option or KINGCLOUD_REGION must be provided in template file or environment variables"))
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}