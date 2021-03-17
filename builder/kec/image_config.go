package kec

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"regexp"
)

type KsyunDiskDevice struct {
	// The instance needs to create a snapshot ID of the image, which must contain a system disk snapshot ID
	// Can be default: Yes, this parameter cannot be default when creating image based on snapshot
	SnapshotId string `mapstructure:"snapshot_id" required:"false"`
	// The ID of the data disk that the instance needs to mirror
	// Default: Yes
	DataDiskId string `mapstructure:"data_disk_id" required:"false"`
}

type KsyunDiskDevices struct {
	SnapshotIds []KsyunDiskDevice `mapstructure:"snapshot_ids" required:"false"`
	DataDiskIds []KsyunDiskDevice `mapstructure:"data_disk_ids" required:"false"`
}

type KsyunImageConfig struct {
	// The name of the user-defined image, [2, 64] English or Chinese
	// characters. It must begin with an uppercase/lowercase letter or a
	// Chinese character, and may contain numbers, `_` or `-`. It cannot begin
	// with `http://` or `https://`.
	KsyunImageName string `mapstructure:"image_name" required:"true"`
	// The type of image
	// LocalImage (ebs) or CommonImage (ks3)
	KsyunImageType string `mapstructure:"image_type" required:"false"`

	KsyunDiskDevices `mapstructure:",squash"`
}

func (c *KsyunImageConfig) Prepare(ctx *interpolate.Context) []error {
	var errs []error
	if c.KsyunImageName == "" {
		errs = append(errs, fmt.Errorf("image_name must be specified"))
	} else if len(c.KsyunImageName) < 2 || len(c.KsyunImageName) > 64 {
		errs = append(errs, fmt.Errorf("image_name must less than 64 letters and more than 1 letters"))
	}
	match, _ := regexp.MatchString("^([\\w-@#.\\p{L}]){2,64}$", c.KsyunImageName)
	if !match {
		errs = append(errs, fmt.Errorf("image_name can't matched"))
	}
	return errs
}
