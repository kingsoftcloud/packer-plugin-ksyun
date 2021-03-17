package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/random"
)

type stepCreateKingcloudKec struct {
	KingcloudRunConfig *KingcloudRunConfig
	InstanceId string
}

func (s *stepCreateKingcloudKec) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)
	chargeTypes := []string{"Daily","HourlyInstantSettlement"}

	if s.KingcloudRunConfig.InstanceName == "" {
		s.KingcloudRunConfig.InstanceName = defaultKecInstanceName
	}

	if s.KingcloudRunConfig.InstanceType == "" {
		s.KingcloudRunConfig.InstanceType = defaultKecInstanceType
	}

	ui.Say("Creating Kingcloud Kec Instance")

	createInstance := make(map[string]interface{})
	createInstance["InstanceType"] = s.KingcloudRunConfig.InstanceType
	createInstance["InstanceName"] = s.KingcloudRunConfig.InstanceName
	createInstance["ImageId"] = s.KingcloudRunConfig.SourceImageId
	createInstance["ProjectId"] = s.KingcloudRunConfig.ProjectId
	createInstance["MaxCount"] = "1"
	createInstance["MinCount"] = "1"
	//SystemDisk
	if s.KingcloudRunConfig.SystemDiskType != ""{
		createInstance["SystemDisk.DiskType"] = s.KingcloudRunConfig.SystemDiskType
	}

	if s.KingcloudRunConfig.SystemDiskSize != 0{
		createInstance["SystemDisk.DiskSize"] = s.KingcloudRunConfig.SystemDiskSize
	}
	//localDataDisk
	if s.KingcloudRunConfig.DataDiskGb != 0{
		createInstance["DataDiskGb"] = s.KingcloudRunConfig.DataDiskGb
	}

	// LocalVolumeSnapshotId
	if s.KingcloudRunConfig.LocalVolumeSnapshotId !=""{
		createInstance["LocalVolumeSnapshotId"] = s.KingcloudRunConfig.LocalVolumeSnapshotId
	}

	//ebsDataDisk
	if len(s.KingcloudRunConfig.KingcloudEbsDataDisks) >0 {
		for index,v:=range s.KingcloudRunConfig.KingcloudEbsDataDisks{
			ebsType := fmt.Sprintf("DataDisk.%d.Type",index+1)
			ebsSize := fmt.Sprintf("DataDisk.%d.Size",index+1)
			ebsSnapshotId := fmt.Sprintf("DataDisk.%d.SnapshotId",index+1)
			ebsDelete := fmt.Sprintf("DataDisk.%d.DeleteWithInstance",index+1)
			createInstance[ebsType]=v.EbsDataDiskType
			createInstance[ebsSize]=v.EbsDataDiskSize
			createInstance[ebsDelete] = true
			if v.EbsDataDiskSnapshotId != ""{
				createInstance[ebsSnapshotId]=v.EbsDataDiskSnapshotId
			}
		}
	}
	//subnetId
	createInstance["SubnetId"] = s.KingcloudRunConfig.SubnetId
	//SecurityGroupId
	createInstance["SecurityGroupId"] = s.KingcloudRunConfig.SecurityGroupId
	//PrivateIpAddress
	if s.KingcloudRunConfig.PrivateIpAddress != ""{
		createInstance["PrivateIpAddress"] = s.KingcloudRunConfig.PrivateIpAddress
	}
	createInstance["KeepImageLogin"] = false

	//password/ssh/key
	if s.KingcloudRunConfig.Comm.SSHAgentAuth {
		createInstance["KeyId.1"] = s.KingcloudRunConfig.Comm.SSHKeyPairName
	}else{
		s.KingcloudRunConfig.Comm.SSHUsername = defaultKecSshUserName
		if s.KingcloudRunConfig.Comm.SSHPassword != "" {
			createInstance["InstancePassword"] = s.KingcloudRunConfig.Comm.SSHPassword
		}else{
			password := random.AlphaNumUpper(4) + random.AlphaNum(4)+random.AlphaNumLower(4)
			s.KingcloudRunConfig.Comm.SSHPassword = password
			createInstance["InstancePassword"] = password
		}
	}
	//chargeType
	checkChargeType := false
	if s.KingcloudRunConfig.InstanceChargeType == ""{
		s.KingcloudRunConfig.InstanceChargeType = defaultKecChargeType
		checkChargeType = true
	}else{
		for _,v:= range chargeTypes{
			if s.KingcloudRunConfig.InstanceChargeType == v{
				checkChargeType = true
				break
			}
		}
	}
	if !checkChargeType{
		return Halt(stateBag,fmt.Errorf("instance_charge_type not match"),"")
	}
	createInstance["ChargeType"] = s.KingcloudRunConfig.InstanceChargeType
	//SriovNetSupport
	if s.KingcloudRunConfig.SriovNetSupport{
		createInstance["SriovNetSupport"] = s.KingcloudRunConfig.SriovNetSupport
	}
	//userdata
	if s.KingcloudRunConfig.UserData != ""{
		createInstance["UserData"] = s.KingcloudRunConfig.UserData
	}
	//create
	createResp,createErr :=client.KecClient.RunInstances(&createInstance)
	if createErr !=nil {
		return Halt(stateBag, createErr, "Error creating new kec instance")
	}
	if createResp!= nil{
		//Get data
		instanceId :=getSdkValue(stateBag,"InstancesSet.0.InstanceId",*createResp).(string)
		s.InstanceId = instanceId
		// wait
		ui.Say("Waiting Kingcloud Kec Instance Active")
		_,waitErr := client.WaitKecInstanceStatus(stateBag,instanceId,s.KingcloudRunConfig.ProjectId,"active")
		if waitErr !=nil {
			return Halt(stateBag, createErr, fmt.Sprintf("Error Wait new kec instance id %s status active",instanceId))
		}
		stateBag.Put("InstanceId", instanceId)
	}
	return multistep.ActionContinue
}

func (s *stepCreateKingcloudKec) Cleanup(stateBag multistep.StateBag) {
	if s.InstanceId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Kec Instance with Id %s ",s.InstanceId))
		deleteInstance := make(map[string]interface{})
		deleteInstance["InstanceId.1"] = s.InstanceId
		deleteInstance["ForceDelete"] = true
		_,err := client.KecClient.TerminateInstances(&deleteInstance)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Kec Instance %s", err))
		}
	}
}



