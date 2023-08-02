package kec

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type stepCreateKsyunKec struct {
	KsyunRunConfig *KsyunKecRunConfig
	InstanceId     string
}

func (s *stepCreateKsyunKec) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientKecWrapper)
	chargeTypes := []string{"Daily", "HourlyInstantSettlement"}

	if s.KsyunRunConfig.InstanceName == "" {
		s.KsyunRunConfig.InstanceName = defaultKecInstanceName
	}

	if s.KsyunRunConfig.InstanceType == "" {
		s.KsyunRunConfig.InstanceType = defaultKecInstanceType
	}

	ui.Say("Creating Ksyun Kec Instance")

	sourceImage, ok := stateBag.Get("source_image").(*ksyun.Ks3Image)
	if !ok {
		stateBag.Put("error", fmt.Errorf("source_image type assertion failed"))
		return multistep.ActionHalt
	}
	s.KsyunRunConfig.SourceImageId = sourceImage.ImageId

	createInstance := make(map[string]interface{})
	createInstance["InstanceType"] = s.KsyunRunConfig.InstanceType
	createInstance["InstanceName"] = s.KsyunRunConfig.InstanceName
	createInstance["ImageId"] = s.KsyunRunConfig.SourceImageId
	createInstance["ProjectId"] = s.KsyunRunConfig.ProjectId
	createInstance["MaxCount"] = "1"
	createInstance["MinCount"] = "1"
	// SystemDisk
	if s.KsyunRunConfig.SystemDiskType != "" {
		createInstance["SystemDisk.DiskType"] = s.KsyunRunConfig.SystemDiskType
	}

	if s.KsyunRunConfig.SystemDiskSize != 0 {
		createInstance["SystemDisk.DiskSize"] = s.KsyunRunConfig.SystemDiskSize
	}
	// localDataDisk
	if s.KsyunRunConfig.DataDiskGb != 0 {
		createInstance["DataDiskGb"] = s.KsyunRunConfig.DataDiskGb
	}

	// LocalVolumeSnapshotId
	if s.KsyunRunConfig.LocalVolumeSnapshotId != "" {
		createInstance["LocalVolumeSnapshotId"] = s.KsyunRunConfig.LocalVolumeSnapshotId
	}

	// ebsDataDisk
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
	// subnetId
	createInstance["subnetId"] = s.KsyunRunConfig.SubnetId
	// securityGroupId
	createInstance["securityGroupId"] = s.KsyunRunConfig.SecurityGroupId
	// PrivateIpAddress
	if s.KsyunRunConfig.PrivateIpAddress != "" {
		createInstance["PrivateIpAddress"] = s.KsyunRunConfig.PrivateIpAddress
	}
	createInstance["KeepImageLogin"] = false

	// password/ssh/key
	if s.KsyunRunConfig.Comm.SSHKeyPairName != "" {
		createInstance["KeyId.1"] = s.KsyunRunConfig.Comm.SSHKeyPairName
	} else {
		if s.KsyunRunConfig.Comm.SSHPassword != "" {
			createInstance["InstancePassword"] = s.KsyunRunConfig.Comm.SSHPassword
		} else {
			createInstance["InstancePassword"] = s.KsyunRunConfig.Comm.WinRMPassword
		}
	}
	// chargeType
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
		return ksyun.Halt(stateBag, fmt.Errorf("instance_charge_type not match"), "")
	}
	createInstance["ChargeType"] = s.KsyunRunConfig.InstanceChargeType
	// SriovNetSupport
	if s.KsyunRunConfig.SriovNetSupport {
		createInstance["SriovNetSupport"] = s.KsyunRunConfig.SriovNetSupport
	}
	// userdata
	if s.KsyunRunConfig.UserData != "" {
		createInstance["UserData"] = s.KsyunRunConfig.UserData
	}
	// create
	createResp, createErr := client.KecClient.RunInstances(&createInstance)
	if createErr != nil {
		return ksyun.Halt(stateBag, createErr, "Error creating new kec instance")
	}
	if createResp != nil {
		// Get data
		instanceId := ksyun.GetSdkValue(stateBag, "InstancesSet.0.InstanceId", *createResp).(string)

		s.InstanceId = instanceId
		// wait
		ui.Say("Waiting Ksyun Kec Instance Active")
		_, waitErr := client.WaitKecInstanceStatus(stateBag, instanceId, s.KsyunRunConfig.ProjectId, "active")
		if waitErr != nil {
			return ksyun.Halt(stateBag, createErr, fmt.Sprintf("Error Wait new kec instance id %s status active", instanceId))
		}
		stateBag.Put("InstanceId", instanceId)

		// processing tag on kec
		if len(s.KsyunRunConfig.RunTags) > 0 {
			ui.Say("Pinning tags on Kec instance")
			ksyunTags := ksyun.TagMap(s.KsyunRunConfig.RunTags).KsyunTags()
			ksyunTags.Report(ui)

			err := ksyunTags.Pinning(ksyun.ResourceTypeKec, instanceId, client.TagsClient)
			if err != nil {
				return ksyun.Halt(stateBag, err, "Error pinning tags to instance")
			}
		}
	}
	return multistep.ActionContinue
}

func (s *stepCreateKsyunKec) Cleanup(stateBag multistep.StateBag) {
	if s.InstanceId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("client").(*ClientKecWrapper)
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
