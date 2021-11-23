package epc

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"regexp"
)

type KsyunImageConfig struct {
	KsyunImageName string `mapstructure:"image_name" required:"false"`
}

func (c *KsyunImageConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	if c.KsyunImageName == "" {
		c.KsyunImageName = defaultEpcImageName
	}
	if len(c.KsyunImageName) < 2 || len(c.KsyunImageName) > 64 {
		errs = append(errs, fmt.Errorf("image_name must less than 64 letters and more than 1 letters"))
	}
	match, _ := regexp.MatchString("^([\\w-@#.\\p{L}]){2,64}$", c.KsyunImageName)
	if !match {
		errs = append(errs, fmt.Errorf("image_name can't matched"))
	}
	return errs
}
