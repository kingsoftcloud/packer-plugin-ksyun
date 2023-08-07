//go:generate packer-sdc struct-markdown
//go:generate packer-sdc mapstructure-to-hcl2 -type KmiFilterOptions

package ksyun

import (
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
)

type CommonConfig struct {
	// VPC ID allocated by the system.
	VpcId string `mapstructure:"vpc_id" required:"false"`
	// the default value is packer_vpc
	VpcName string `mapstructure:"vpc_name" required:"false"`
	// 172.16.0.0/16. When not specified, the default value is 172.16.0.0/16.
	VpcCidrBlock string `mapstructure:"vpc_cidr_block" required:"false"`
	// The ID of the Subnet to be used.
	SubnetId string `mapstructure:"subnet_id" required:"false"`
	// the default value is packer_subnet
	SubnetName string `mapstructure:"subnet_name" required:"false"`
	// 172.16.0.0/24. When not specified, the default value is 172.16.0.0/24.
	DNS1            string `mapstructure:"dns1" required:"false"`
	DNS2            string `mapstructure:"dns2" required:"false"`
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
	// the default value is packer_security_group
	SecurityGroupName string `mapstructure:"security_group_name" required:"false"`
	// Private IP address, which specifies any valid value within the range of subnet IP address and represents
	// the primary IP address of the instance. Only one can be selected and bound to the primary network card.
	// If this parameter is not specified, the system will automatically select one from the valid address pool at random
	// Valid values: standard IP address format
	PrivateIpAddress string `mapstructure:"private_ip_address" required:"false"`
	// Indicating associate whether public ip for creating kec instance
	// Default `false`
	AssociatePublicIpAddress bool `mapstructure:"associate_public_ip_address" required:"false"`
	// PublicIp charge type, which can be
	// Daily TrafficMonthly DailyPaidByTransfer HourlyInstantSettlement
	// Default is Daily
	PublicIpChargeType string `mapstructure:"public_ip_charge_type" required:"false"`
	// [1-100]
	// Default is 1
	PublicIpBandWidth int `mapstructure:"public_ip_band_width" required:"false"`
	// Default is 0
	ProjectId string `mapstructure:"project_id" required:"false"`

	// Key/value pair tags to apply to the instance that is launched to create the image.
	// These tags are not applied to the resulting image.
	RunTags map[string]string `mapstructure:"run_tags" required:"false"`
	// Same as [`run_tags`](#run_tags) but defined as a singular repeatable
	// block containing a `key` and a `value` field. In HCL2 mode the
	// [`dynamic_block`](/packer/docs/templates/hcl_templates/expressions#dynamic-blocks)
	// will allow you to create those programatically.
	RunTag config.KeyValues `mapstructure:"run_tag" required:"false"`

	// Communicator settings
	Comm communicator.Config `mapstructure:",squash"`
}

type KmiFilterOptions struct {
	// Selects the newest created image when true.
	// This is most useful for selecting a daily distro build.
	MostRecent bool `mapstructure:"most_recent"`

	// ImageSource Valid values are import, copy, share, extend, system.
	ImageSource string `mapstructure:"image_source"`

	// NameRegex A regex string to filter resulting images by name.
	// (Such as: `^CentOS 7.[1-2] 64` means CentOS 7.1 of 64-bit operating system or CentOS 7.2 of 64-bit operating system,
	// \"^Ubuntu 16.04 64\" means Ubuntu 16.04 of 64-bit operating system).
	NameRegex string `mapstructure:"name_regex"`

	// Platform type of the image system.
	Platform string `mapstructure:"platform"`
}
