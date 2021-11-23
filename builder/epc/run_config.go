package epc

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	"github.com/hashicorp/packer-plugin-sdk/uuid"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type KsyunEpcRunConfig struct {
	HostType                   string `mapstructure:"host_type" required:"true"`
	Raid                       string `mapstructure:"raid" required:"false"`
	SourceImageId              string `mapstructure:"source_image_id" required:"true"`
	NetworkInterfaceMode       string `mapstructure:"network_interface_mode" required:"false"`
	HostName                   string `mapstructure:"host_name" required:"false"`
	ComputerName               string `mapstructure:"computer_name" required:"false"`
	HostChargeType             string `mapstructure:"host_charge_type" required:"false"`
	SecurityAgent              string `mapstructure:"security_agent" required:"false"`
	ContainerAgent             string `mapstructure:"container_agent" required:"false"`
	CloudMonitorAgent          string `mapstructure:"cloud_monitor_agent" required:"false"`
	SystemFileType             string `mapstructure:"system_file_type" required:"false"`
	DataFileType               string `mapstructure:"data_file_type" required:"false"`
	DataDiskCatalogue          string `mapstructure:"data_disk_catalogue" required:"false"`
	DataDiskCatalogueSuffix    string `mapstructure:"data_disk_catalogue_suffix" required:"false"`
	ExtensionSubnetId          string `mapstructure:"extension_subnet_id" required:"false"`
	ExtensionSubnetName        string `mapstructure:"extension_subnet_name" required:"false"`
	ExtensionSubnetCidrBlock   string `mapstructure:"extension_subnet_cidr_block" required:"false"`
	ExtensionPrivateIpAddress  string `mapstructure:"extension_private_ip_address" required:"false"`
	ExtensionDNS1              string `mapstructure:"extension_dns1" required:"false"`
	ExtensionDNS2              string `mapstructure:"extension_dns2" required:"false"`
	ExtensionSecurityGroupId   string `mapstructure:"extension_security_group_id" required:"false"`
	ExtensionSecurityGroupName string `mapstructure:"extension_security_group_name" required:"false"`
	TempSubnetId               string
	TempSecurityGroupId        string
	ksyun.CommonConfig         `mapstructure:",squash"`
}

func (c *KsyunEpcRunConfig) Prepare(ctx *interpolate.Context) []error {
	errs := c.Init(ctx)
	return append(errs, c.Check()...)
}

func (c *KsyunEpcRunConfig) Init(ctx *interpolate.Context) []error {
	if c.Comm.SSHKeyPairName == "" && c.Comm.SSHTemporaryKeyPairName == "" &&
		c.Comm.SSHPrivateKeyFile == "" && c.Comm.SSHPassword == "" && c.Comm.WinRMPassword == "" {
		c.Comm.SSHTemporaryKeyPairName = fmt.Sprintf("packer_epc_%s", uuid.TimeOrderedUUID())
	}
	if c.Raid == "" {
		c.Raid = "SRaid0"
	}
	if c.NetworkInterfaceMode == "" {
		c.NetworkInterfaceMode = "bond4"
	}
	if c.HostName == "" {
		c.HostName = defaultEpcInstanceName
	}
	if c.HostChargeType == "" {
		c.HostChargeType = defaultEpcChargeType
	}
	if c.SecurityAgent == "" {
		c.SecurityAgent = "no"
	}
	if c.CloudMonitorAgent == "" {
		c.CloudMonitorAgent = "no"
	}
	if c.ContainerAgent == "" {
		c.ContainerAgent = "unsupported"
	}
	if c.SystemFileType == "" {
		c.SystemFileType = "EXT4"
	}
	if c.DataFileType == "" {
		c.DataFileType = "XFS"
	}
	if c.DataDiskCatalogue == "" {
		c.DataDiskCatalogue = "/DATA/disk"
	}
	if c.DataDiskCatalogueSuffix == "" {
		c.DataDiskCatalogueSuffix = "NaturalNumber"
	}
	return c.Comm.Prepare(ctx)
}

func (c *KsyunEpcRunConfig) Check() (errors []error) {
	if !ksyun.StringInSlice(c.Raid, []string{"Raid1", "Raid5", "Raid10", "Raid50", "SRaid0"}, false) {
		errors = append(errors, fmt.Errorf(" Raid is inValid"))
	}
	if !ksyun.StringInSlice(c.NetworkInterfaceMode, []string{"bond4", "single", "dual"}, false) {
		errors = append(errors, fmt.Errorf("NetworkInterfaceMode is inValid"))
	}
	if !ksyun.StringInSlice(c.HostChargeType, []string{"Daily"}, false) {
		errors = append(errors, fmt.Errorf("HostChargeType is inValid"))
	}
	if !ksyun.StringInSlice(c.SecurityAgent, []string{"classic", "no"}, false) {
		errors = append(errors, fmt.Errorf("SecurityAgent is inValid"))
	}
	if !ksyun.StringInSlice(c.CloudMonitorAgent, []string{"classic", "no"}, false) {
		errors = append(errors, fmt.Errorf("CloudMonitorAgent is inValid"))
	}
	if !ksyun.StringInSlice(c.ContainerAgent, []string{"supported", "unsupported"}, false) {
		errors = append(errors, fmt.Errorf("ContainerAgent is inValid"))
	}
	if !ksyun.StringInSlice(c.SystemFileType, []string{"EXT4", "XFS"}, false) {
		errors = append(errors, fmt.Errorf("SystemFileType is inValid"))
	}
	if !ksyun.StringInSlice(c.DataFileType, []string{"EXT4", "XFS"}, false) {
		errors = append(errors, fmt.Errorf("DataFileType is inValid"))
	}
	if !ksyun.StringInSlice(c.DataDiskCatalogue, []string{"/DATA/disk", "/data"}, false) {
		errors = append(errors, fmt.Errorf("DataDiskCatalogue is inValid"))
	}
	if !ksyun.StringInSlice(c.DataDiskCatalogueSuffix, []string{"NoSuffix", "NaturalNumber", "NaturalNumberFromZero"}, false) {
		errors = append(errors, fmt.Errorf("DataDiskCatalogueSuffix is inValid"))
	}
	return errors
}

func (c *KsyunEpcRunConfig) NeedExtensionNetwork() bool {
	if c.NetworkInterfaceMode == "dual" {
		return true
	}
	return false
}

func (c *KsyunEpcRunConfig) ExtensionSubnet(config *ksyun.CommonConfig) ksyun.AfterStepRun {
	return func() {

	}
}

func (c *KsyunEpcRunConfig) PrepareExtensionSubnet(config *ksyun.CommonConfig) ksyun.AfterStepRun {
	return func() {
		if c.NeedExtensionNetwork() {
			c.TempSubnetId = config.SubnetId
			config.SubnetId = c.ExtensionSubnetId
			config.SubnetName = c.ExtensionSubnetName
			config.SubnetCidrBlock = c.ExtensionSubnetCidrBlock
			config.PrivateIpAddress = c.ExtensionPrivateIpAddress
			config.DNS1 = c.ExtensionDNS1
			config.DNS2 = c.ExtensionDNS2
		}
	}
}

func (c *KsyunEpcRunConfig) MergeExtensionSubnet(config *ksyun.CommonConfig) ksyun.AfterStepRun {
	return func() {
		if c.NeedExtensionNetwork() {
			c.ExtensionSubnetId = config.SubnetId
			config.SubnetId = c.TempSubnetId
		}
	}
}

func (c *KsyunEpcRunConfig) PrepareExtensionSecurityGroup(config *ksyun.CommonConfig) ksyun.AfterStepRun {
	return func() {
		if c.NeedExtensionNetwork() {
			c.TempSecurityGroupId = config.SecurityGroupId
			config.SecurityGroupId = c.ExtensionSecurityGroupId
			config.SecurityGroupName = c.ExtensionSecurityGroupName
		}
	}
}

func (c *KsyunEpcRunConfig) MergeExtensionSecurityGroup(config *ksyun.CommonConfig) ksyun.AfterStepRun {
	return func() {
		if c.NeedExtensionNetwork() {
			c.ExtensionSecurityGroupId = config.SecurityGroupId
			config.SecurityGroupId = c.TempSecurityGroupId
		}
	}
}
