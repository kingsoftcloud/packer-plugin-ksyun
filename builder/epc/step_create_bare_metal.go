package epc

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type stepCreateBareMetal struct {
	KsyunRunConfig *KsyunEpcRunConfig
	bareMetalId    string
}

func (s *stepCreateBareMetal) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientEpcWrapper)
	ui.Say("Creating Ksyun Bare Metal Instance")
	createInstance := make(map[string]interface{})
	createInstance["AvailabilityZone"] = s.KsyunRunConfig.AvailabilityZone
	createInstance["HostType"] = s.KsyunRunConfig.HostType
	createInstance["Raid"] = s.KsyunRunConfig.Raid
	createInstance["ImageId"] = s.KsyunRunConfig.SourceImageId
	createInstance["HostName"] = s.KsyunRunConfig.HostName
	createInstance["ChargeType"] = s.KsyunRunConfig.HostChargeType
	createInstance["KeyId"] = s.KsyunRunConfig.Comm.SSHKeyPairName
	createInstance["SecurityAgent"] = s.KsyunRunConfig.SecurityAgent
	createInstance["CloudMonitorAgent"] = s.KsyunRunConfig.CloudMonitorAgent
	createInstance["ContainerAgent"] = s.KsyunRunConfig.ContainerAgent
	createInstance["NetworkInterfaceMode"] = s.KsyunRunConfig.NetworkInterfaceMode
	createInstance["SystemFileType"] = s.KsyunRunConfig.SystemFileType
	createInstance["DataFileType"] = s.KsyunRunConfig.DataFileType
	createInstance["DataDiskCatalogue"] = s.KsyunRunConfig.DataDiskCatalogue
	createInstance["DataDiskCatalogueSuffix"] = s.KsyunRunConfig.DataDiskCatalogueSuffix
	createInstance["SubnetId"] = s.KsyunRunConfig.SubnetId
	createInstance["SecurityGroupId.1"] = s.KsyunRunConfig.SecurityGroupId
	if s.KsyunRunConfig.PrivateIpAddress != "" {
		createInstance["PrivateIpAddress"] = s.KsyunRunConfig.PrivateIpAddress
	}
	if s.KsyunRunConfig.DNS1 != "" {
		createInstance["DNS1"] = s.KsyunRunConfig.DNS1
	}
	if s.KsyunRunConfig.DNS2 != "" {
		createInstance["DNS2"] = s.KsyunRunConfig.DNS2
	}
	if s.KsyunRunConfig.ComputerName != "" {
		createInstance["ComputerName"] = s.KsyunRunConfig.ComputerName
	}
	if s.KsyunRunConfig.NeedExtensionNetwork() {
		createInstance["ExtensionSubnetId"] = s.KsyunRunConfig.ExtensionSubnetId
		createInstance["ExtensionSecurityGroupId.1"] = s.KsyunRunConfig.ExtensionSecurityGroupId
		if s.KsyunRunConfig.ExtensionPrivateIpAddress != "" {
			createInstance["ExtensionPrivateIpAddress"] = s.KsyunRunConfig.ExtensionPrivateIpAddress
		}
		if s.KsyunRunConfig.ExtensionDNS1 != "" {
			createInstance["ExtensionDNS1"] = s.KsyunRunConfig.ExtensionDNS1
		}
		if s.KsyunRunConfig.ExtensionDNS2 != "" {
			createInstance["ExtensionDNS2"] = s.KsyunRunConfig.ExtensionDNS2
		}
	}

	createResp, createErr := client.EpcClient.CreateEpc(&createInstance)
	if createErr != nil {
		return ksyun.Halt(stateBag, createErr, "Error creating new Bare metal instance")
	}
	if createResp != nil {
		//Get data
		instanceId := ksyun.GetSdkValue(stateBag, "Host.HostId", *createResp).(string)
		s.bareMetalId = instanceId
		// wait
		ui.Say(fmt.Sprintf("Waiting Ksyun Bare Metal Instance instance id %s Running", s.bareMetalId))
		_, waitErr := client.WaitEpcInstanceStatus(stateBag, instanceId, s.KsyunRunConfig.ProjectId, "Running")
		if waitErr != nil {
			return ksyun.Halt(stateBag, createErr, fmt.Sprintf("Error Wait new Bare metal instance id %s status Running", instanceId))
		}
		stateBag.Put("InstanceId", instanceId)
	}
	return multistep.ActionContinue
}

func (s *stepCreateBareMetal) Cleanup(stateBag multistep.StateBag) {

	if s.bareMetalId != "" {
		client := stateBag.Get("client").(*ClientEpcWrapper)
		ui := stateBag.Get("ui").(packersdk.Ui)
		ui.Say(fmt.Sprintf("Deleting Bare Metal Instance with Id %s ", s.bareMetalId))
		deleteInstance := make(map[string]interface{})
		deleteInstance["HostId"] = s.bareMetalId
		_, err := client.EpcClient.DeleteEpc(&deleteInstance)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Bare Metal Instance %s", err))
		}
	}
}
