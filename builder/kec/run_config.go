//go:generate struct-markdown

package kec

import (
	"errors"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
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

type KsyunRunConfig struct {
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
	// VPC ID allocated by the system.
	VpcId string `mapstructure:"vpc_id" required:"false"`
	// The VPC name. The default value is blank. [2, 128]
	// English or Chinese characters, must begin with an uppercase/lowercase
	// letter or Chinese character. Can contain numbers, _ and -. The disk
	// description will appear on the console. Cannot begin with `http://` or
	// `https://`.
	//the default value is packer_vpc
	VpcName string `mapstructure:"vpc_name" required:"false"`
	// 172.16.0.0/16. When not specified, the default value is 172.16.0.0/16.
	VpcCidrBlock string `mapstructure:"vpc_cidr_block" required:"false"`
	// The ID of the Subnet to be used.
	SubnetId string `mapstructure:"subnet_id" required:"false"`
	//the default value is packer_subnet
	SubnetName string `mapstructure:"subnet_name" required:"false"`
	// 172.16.0.0/24. When not specified, the default value is 172.16.0.0/24.
	SubnetCidrBlock string `mapstructure:"subnet_cidr_block" required:"false"`
	// availability_zone
	AvailabilityZone string `mapstructure:"availability_zone" required:"false"`
	// ID of the security group to which a newly
	// created instance belongs. Mutual access is allowed between instances in one
	// security group. If not specified, the newly created instance will be added
	// to the default security group. If the default group doesn’t exist, or the
	// number of instances in it has reached the maximum limit, a new security
	// group will be created automatically.
	SecurityGroupId string `mapstructure:"security_group_id" required:"false"`
	// The security group name. The default value
	// is blank. [2, 128] English or Chinese characters, must begin with an
	// uppercase/lowercase letter or Chinese character. Can contain numbers, .,
	// _ or -. It cannot begin with `http://` or `https://`.
	//the default value is packer_security_group
	SecurityGroupName string `mapstructure:"security_group_name" required:"false"`
	// Private IP address, which specifies any valid value within the range of subnet IP address and represents
	// the primary IP address of the instance. Only one can be selected and bound to the primary network card.
	// If this parameter is not specified, the system will automatically select one from the valid address pool at random
	// Valid values: standard IP address format
	PrivateIpAddress string `mapstructure:"private_ip_address" required:"false"`
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
	SriovNetSupport bool `mapstructure:"sriov_net_support" required:"false"`
	// Default is 0
	ProjectId                string `mapstructure:"project_id" required:"false"`
	AssociatePublicIpAddress bool   `mapstructure:"associate_public_ip_address" required:"false"`
	// PublicIp charge type, which can be
	// Daily TrafficMonthly DailyPaidByTransfer HourlyInstantSettlement
	// Default is Daily
	PublicIpChargeType string `mapstructure:"public_ip_charge_type" required:"false"`
	// [1-100]
	// Default is 1
	PublicIpBandWidth int `mapstructure:"public_ip_band_width" required:"false"`
	// The local data disk snapshot ID can be used to create a data disk based on the snapshot.
	// This parameter is valid only if datadiskgb is specified and the size is the same as the snapshot size
	LocalVolumeSnapshotId string `mapstructure:"local_volume_snapshot_id" required:"false"`
	// The data volume capacity is in GB. The capacity limit changes according to the definition of the instance package
	// type. If it is not specified when calling, it is the minimum value of the corresponding instance package type.
	// When the instancetype is a general-purpose host N1, N2, N3 or a local nvme I4, this parameter
	// does not take effect.
	DataDiskGb int `mapstructure:"data_disk_gb" required:"false"`
	// The user-defined data provided for instance startup needs to be encoded in Base64 mode,
	// and the maximum data size supported is 16kb
	UserData string `mapstructure:"user_data" required:"false"`
	// Communicator settings
	Comm communicator.Config `mapstructure:",squash"`
}

func (c *KsyunRunConfig) Prepare(ctx *interpolate.Context) []error {
	// SSH Validation
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
		errs = append(errs, errors.New("A kingcloud_instance_type must be specified"))
	}

	return errs
}
