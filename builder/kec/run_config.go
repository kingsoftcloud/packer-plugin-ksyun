//go:generate packer-sdc struct-markdown
package kec

import (
	"errors"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
	"regexp"
)

type KsyunEbsDataDisk struct {
	//  SSD3.0|| EHDD
	EbsDataDiskType string `mapstructure:"data_disk_type" required:"true"`
	//  [10，16000]
	EbsDataDiskSize int `mapstructure:"data_disk_size" required:"true"`
	// pattern ^[a-zA-Z0-9-]{36}$
	EbsDataDiskSnapshotId string `mapstructure:"data_disk_snapshot_id" required:"false"`
}

type KsyunKecRunConfig struct {
	//Instance package type, if the instance package type is not specified when calling, the default value is I1.1A.
	InstanceType  string `mapstructure:"instance_type" required:"true"`
	SourceImageId string `mapstructure:"source_image_id" required:"true"`
	// Local_SSD || SSD3.0 || EHDD
	SystemDiskType string `mapstructure:"system_disk_type" required:"false"`
	SystemDiskSize int    `mapstructure:"system_disk_size" required:"false"`
	// EbsDataDisk
	KsyunEbsDataDisks []KsyunEbsDataDisk `mapstructure:"data_disks" required:"false"`
	// PostPaidByDay or PostPaidByHour
	// default is PostPaidByDay
	InstanceChargeType string `mapstructure:"instance_charge_type" required:"true"`
	// Display name of the instance, which is a string of 2 to 128 Chinese or
	// English characters. It must begin with an uppercase/lowercase letter or
	// a Chinese charac displayed on the Alibaba Cloud console. If this
	//	// parameter is not specified, the default value is InstanceId of the
	//	// instance. It cannot begin with `http://` or `https://`.ter and can contain numerals, `.`, `_`, or `-`. The
	// instance name is packer_kec_instance
	InstanceName string `mapstructure:"instance_name" required:"false"`
	// This parameter needs to satisfy the following two conditions:
	// IO optimized I1, calculation optimized C1, and IO optimized I2 are more than 8C packages
	// We use the standard image improved by Jinshan cloud or the user-defined image made by
	// the instance of starting the Jinshan cloud standard image
	// default : false
	SriovNetSupport       bool   `mapstructure:"sriov_net_support" required:"false"`
	LocalVolumeSnapshotId string `mapstructure:"local_volume_snapshot_id" required:"false"`
	// The data volume capacity is in GB. The capacity limit changes according to the definition of the instance package
	// type. If it is not specified when calling, it is the minimum value of the corresponding instance package type.
	// When the instancetype is a general-purpose host N1, N2, N3 or a local nvme I4, this parameter
	// does not take effect.
	DataDiskGb int `mapstructure:"data_disk_gb" required:"false"`
	// The user-defined data provided for instance startup needs to be encoded in Base64 mode,
	// and the maximum data size supported is 16kb
	UserData string `mapstructure:"user_data" required:"false"`

	ksyun.CommonConfig `mapstructure:",squash"`
}

func (c *KsyunKecRunConfig) Prepare(ctx *interpolate.Context) []error {
	if c.Comm.SSHKeyPairName == "" && c.Comm.SSHTemporaryKeyPairName == "" &&
		c.Comm.SSHPrivateKeyFile == "" && c.Comm.SSHPassword == "" && c.Comm.WinRMPassword == "" {
		c.Comm.SSHTemporaryKeyPairName = fmt.Sprintf("packer_%s", uuid.TimeOrderedUUID())
	}

	// Validation
	errs := c.Comm.Prepare(ctx)
	// source_image
	if c.SourceImageId == "" {
		errs = append(errs, errors.New("A source_image_id must be specified"))
	}

	match, _ := regexp.MatchString("^(IMG-)?[a-zA-Z0-9-]{36}$", c.SourceImageId)
	if !match {
		errs = append(errs, fmt.Errorf("source_image_id can't matched"))
	}

	if c.InstanceType == "" {
		errs = append(errs, errors.New("A ksyun_instance_type must be specified"))
	}

	return errs
}
