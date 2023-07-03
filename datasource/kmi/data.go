package kmi

//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type DatasourceOutput,Config
import (
	"fmt"

	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
	"github.com/zclconf/go-cty/cty"
)

type Datasource struct {
	config Config
}

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	DatasourceClientConfig `mapstructure:",squash"`
	ksyun.KmiFilterOptions `mapstructure:",squash"`

	ctx interpolate.Context
}

type DatasourceClientConfig struct {
	ksyun.ClientConfig `mapstructure:",squash"`
	kecClient          *kec.Kec
}

func (c *DatasourceClientConfig) KecClient() *kec.Kec {
	if c.kecClient != nil {
		return c.kecClient
	}

	c.kecClient = kec.SdkNew(ksc.NewClient(c.KsyunAccessKey, c.KsyunSecretKey),
		&ksc.Config{Region: &c.KsyunRegion},
		&utils.UrlInfo{
			UseSSL: true,
		})
	return c.kecClient
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	var errs *packersdk.MultiError
	errs = packersdk.MultiErrorAppend(errs, d.config.ClientConfig.Prepare(&d.config.ctx)...)

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}
	return nil
}

type DatasourceOutput struct {
	// The ID of the image.
	ID string `mapstructure:"id"`
	// The name of the image.
	Name string `mapstructure:"name"`
	// The date of creation of the image.
	CreationDate string `mapstructure:"creation_date"`

	Platform string `mapstructure:"platform"`

	ImageSource string `mapstructure:"image_source"`
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	kecClient := d.config.KecClient()
	if kecClient == nil {
		return cty.Value{}, fmt.Errorf("the current client is nil")
	}
	var params *map[string]interface{}

	image, err := d.config.GetFilteredImage(params, kecClient)
	if err != nil {
		return cty.NullVal(cty.EmptyObject), err
	}

	output := DatasourceOutput{
		ID:           image.ImageId,
		Name:         image.Name,
		CreationDate: image.CreationDate,
		Platform:     image.Platform,
		ImageSource:  image.ImageSource,
	}
	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
