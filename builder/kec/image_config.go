//go:generate packer-sdc struct-markdown

package kec

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type KsyunKecDiskDevice struct {
	// The instance needs to create a snapshot ID of the image, which must contain a system disk snapshot ID
	// Can be default: Yes, this parameter cannot be default when creating image based on snapshot
	SnapshotId string `mapstructure:"snapshot_id" required:"false"`
	// The ID of the data disk that the instance needs to mirror
	// Default: Yes
	DataDiskId string `mapstructure:"data_disk_id" required:"false"`
}

type KsyunKecDiskDevices struct {
	SnapshotIds []KsyunKecDiskDevice `mapstructure:"snapshot_ids" required:"false"`
	DataDiskIds []KsyunKecDiskDevice `mapstructure:"data_disk_ids" required:"false"`
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

	// If this value is true, the image created will not include any snapshot
	// of data disks. The default value is false.
	KsyunImageIgnoreDataDisks bool `mapstructure:"image_ignore_data_disks" required:"false"`

	KsyunKecDiskDevices `mapstructure:",squash"`

	// Copy to the regions.
	KsyunImageCopyRegions []string `mapstructure:"image_copy_regions" required:"false"`
	// The image name in copied regions
	KsyunImageCopyNames []string `mapstructure:"image_copy_names" required:"false"`

	// Share image to other accounts
	KsyunImageShareAccounts []string `mapstructure:"image_share_accounts" required:"false"`

	// Set the image as warm-up for fast boot
	KsyunImageWarmUp bool `mapstructure:"image_warm_up" required:"false"`
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
