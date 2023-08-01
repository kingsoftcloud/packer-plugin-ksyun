//go:generate packer-sdc struct-markdown
package epc

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type KsyunImageConfig struct {
	// The name of the user-defined image, [2, 64] English or Chinese
	// characters. It must begin with an uppercase/lowercase letter or a
	// Chinese character, and may contain numbers, `_` or `-`. It cannot begin
	// with `http://` or `https://`.
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
