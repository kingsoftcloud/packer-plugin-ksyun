package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/random"
)

type stepCreateKsyunKec struct {
	KsyunRunConfig *KsyunRunConfig
	InstanceId     string
}

func (s *stepCreateKsyunKec) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)
	chargeTypes := []string{"Daily", "HourlyInstantSettlement"}

	if s.KsyunRunConfig.InstanceName == "" {
		s.KsyunRunConfig.InstanceName = defaultKecInstanceName
	}

	if s.KsyunRunConfig.InstanceType == "" {
		s.KsyunRunConfig.InstanceType = defaultKecInstanceType
	}

	ui.Say("Creating Ksyun Kec Instance")

	createInstance := make(map[string]interface{})
	createInstance["InstanceType"] = s.KsyunRunConfig.InstanceType
	createInstance["InstanceName"] = s.KsyunRunConfig.InstanceName
	createInstance["ImageId"] = s.KsyunRunConfig.SourceImageId
	createInstance["ProjectId"] = s.KsyunRunConfig.ProjectId
	createInstance["MaxCount"] = "1"
	createInstance["MinCount"] = "1"
	//SystemDisk
	if s.KsyunRunConfig.SystemDiskType != "" {
		createInstance["SystemDisk.DiskType"] = s.KsyunRunConfig.SystemDiskType
	}

	if s.KsyunRunConfig.SystemDiskSize != 0 {
		createInstance["SystemDisk.DiskSize"] = s.KsyunRunConfig.SystemDiskSize
	}
	//localDataDisk
	if s.KsyunRunConfig.DataDiskGb != 0 {
		createInstance["DataDiskGb"] = s.KsyunRunConfig.DataDiskGb
	}

	// LocalVolumeSnapshotId
	if s.KsyunRunConfig.LocalVolumeSnapshotId != "" {
		createInstance["LocalVolumeSnapshotId"] = s.KsyunRunConfig.LocalVolumeSnapshotId
	}

	//ebsDataDisk
	if len(s.KsyunRunConfig.KsyunEbsDataDisks) > 0 {
		for index, v := range s.KsyunRunConfig.KsyunEbsDataDisks {
			ebsType := fmt.Sprintf("DataDisk.%d.Type", index+1)
			ebsSize := fmt.Sprintf("DataDisk.%d.Size", index+1)
			ebsSnapshotId := fmt.Sprintf("DataDisk.%d.SnapshotId", index+1)
			ebsDelete := fmt.Sprintf("DataDisk.%d.DeleteWithInstance", index+1)
			createInstance[ebsType] = v.EbsDataDiskType
			createInstance[ebsSize] = v.EbsDataDiskSize
			createInstance[ebsDelete] = true
			if v.EbsDataDiskSnapshotId != "" {
				createInstance[ebsSnapshotId] = v.EbsDataDiskSnapshotId
			}
		}
	}
	//subnetId
	createInstance["SubnetId"] = s.KsyunRunConfig.SubnetId
	//SecurityGroupId
	createInstance["SecurityGroupId"] = s.KsyunRunConfig.SecurityGroupId
	//PrivateIpAddress
	if s.KsyunRunConfig.PrivateIpAddress != "" {
		createInstance["PrivateIpAddress"] = s.KsyunRunConfig.PrivateIpAddress
	}
	createInstance["KeepImageLogin"] = false

	//password/ssh/key
	if s.KsyunRunConfig.Comm.SSHAgentAuth {
		createInstance["KeyId.1"] = s.KsyunRunConfig.Comm.SSHKeyPairName
	} else {
		s.KsyunRunConfig.Comm.SSHUsername = defaultKecSshUserName
		if s.KsyunRunConfig.Comm.SSHPassword != "" {
			createInstance["InstancePassword"] = s.KsyunRunConfig.Comm.SSHPassword
		} else {
			password := random.AlphaNumUpper(4) + random.AlphaNum(4) + random.AlphaNumLower(4)
			s.KsyunRunConfig.Comm.SSHPassword = password
			createInstance["InstancePassword"] = password
		}
	}
	//chargeType
	checkChargeType := false
	if s.KsyunRunConfig.InstanceChargeType == "" {
		s.KsyunRunConfig.InstanceChargeType = defaultKecChargeType
		checkChargeType = true
	} else {
		for _, v := range chargeTypes {
			if s.KsyunRunConfig.InstanceChargeType == v {
				checkChargeType = true
				break
			}
		}
	}
	if !checkChargeType {
		return Halt(stateBag, fmt.Errorf("instance_charge_type not match"), "")
	}
	createInstance["ChargeType"] = s.KsyunRunConfig.InstanceChargeType
	//SriovNetSupport
	if s.KsyunRunConfig.SriovNetSupport {
		createInstance["SriovNetSupport"] = s.KsyunRunConfig.SriovNetSupport
	}
	//userdata
	if s.KsyunRunConfig.UserData != "" {
		createInstance["UserData"] = s.KsyunRunConfig.UserData
	}
	//create
	createResp, createErr := client.KecClient.RunInstances(&createInstance)
	if createErr != nil {
		return Halt(stateBag, createErr, "Error creating new kec instance")
	}
	if createResp != nil {
		//Get data
		instanceId := getSdkValue(stateBag, "InstancesSet.0.InstanceId", *createResp).(string)
		s.InstanceId = instanceId
		// wait
		ui.Say("Waiting Ksyun Kec Instance Active")
		_, waitErr := client.WaitKecInstanceStatus(stateBag, instanceId, s.KsyunRunConfig.ProjectId, "active")
		if waitErr != nil {
			return Halt(stateBag, createErr, fmt.Sprintf("Error Wait new kec instance id %s status active", instanceId))
		}
		stateBag.Put("InstanceId", instanceId)
	}
	return multistep.ActionContinue
}

func (s *stepCreateKsyunKec) Cleanup(stateBag multistep.StateBag) {
	if s.InstanceId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Kec Instance with Id %s ", s.InstanceId))
		deleteInstance := make(map[string]interface{})
		deleteInstance["InstanceId.1"] = s.InstanceId
		deleteInstance["ForceDelete"] = true
		_, err := client.KecClient.TerminateInstances(&deleteInstance)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Kec Instance %s", err))
		}
	}
}
