package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepConfigKingcloudSubnet struct {
	KingcloudRunConfig *KingcloudRunConfig
	subnetId string
}

func (s *stepConfigKingcloudSubnet) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)

	if s.KingcloudRunConfig.SubnetId != ""{
		//Check_Subnet
		querySubnet := make(map[string]interface{})
		querySubnet["SubnetId.1"] = s.KingcloudRunConfig.SubnetId
		resp,err := client.VpcClient.DescribeSubnets(&querySubnet)
		if err != nil {
			return Halt(stateBag, err, fmt.Sprintf("Error query Subnet with id %s",s.KingcloudRunConfig.SubnetId))
		}
		if resp !=nil{
			subnetId := getSdkValue(stateBag,"SubnetSet.0.SubnetId",*resp)
			subnetType := getSdkValue(stateBag,"SubnetSet.0.SubnetType",*resp)
			subnetName := getSdkValue(stateBag,"SubnetSet.0.SubnetName",*resp)
			vpcId := getSdkValue(stateBag,"SubnetSet.0.VpcId",*resp)
			if subnetId == nil {
				return Halt(stateBag, err, fmt.Sprintf("Subnet id %s not found",s.KingcloudRunConfig.SubnetId))
			}

			if vpcId != s.KingcloudRunConfig.VpcId {
				return Halt(stateBag, fmt.Errorf(fmt.Sprintf("Subnet id %s vpc not match",
					s.KingcloudRunConfig.SubnetId)),"" )
			}

			if subnetType != EnableSubnetType {
				return Halt(stateBag,
					fmt.Errorf(fmt.Sprintf("Subnet id %s Type is Not %s",EnableSubnetType,s.KingcloudRunConfig.SubnetId)), "")
			}
			ui.Say(fmt.Sprintf("Using existing Subnet id is %s name is %s", s.KingcloudRunConfig.SubnetId,
				subnetName))
		}
		return multistep.ActionContinue
	}else{
		//create_subnet
		if s.KingcloudRunConfig.SubnetName == ""{
			s.KingcloudRunConfig.SubnetName = defaultSubnetName
		}
		if s.KingcloudRunConfig.SubnetCidrBlock == ""{
			s.KingcloudRunConfig.SubnetCidrBlock = defaultSubnetCidr
		}
		startIp,minIp,maxIp := getCidrIpRange(s.KingcloudRunConfig.SubnetCidrBlock)
		ui.Say(fmt.Sprintf("Creating new Subnet with name  %s cidr %s vpcId %s",
			s.KingcloudRunConfig.SubnetName,s.KingcloudRunConfig.SubnetCidrBlock,s.KingcloudRunConfig.VpcId))
		createSubnet := make(map[string]interface{})
		createSubnet["VpcId"] = s.KingcloudRunConfig.VpcId
		createSubnet["SubnetName"] = s.KingcloudRunConfig.SubnetName
		createSubnet["SubnetType"] = EnableSubnetType
		createSubnet["CidrBlock"] = s.KingcloudRunConfig.SubnetCidrBlock
		createSubnet["GatewayIp"] = startIp
		createSubnet["DhcpIpFrom"] = minIp
		createSubnet["DhcpIpTo"] = maxIp
		if s.KingcloudRunConfig.AvailabilityZone != ""{
			createSubnet["AvailabilityZone"] = s.KingcloudRunConfig.AvailabilityZone
		}
		resp,err := client.VpcClient.CreateSubnet(&createSubnet)
		if err != nil {
			return Halt(stateBag, err, "Error creating new Subnet")
		}
		if resp !=nil {
			s.KingcloudRunConfig.SubnetId = getSdkValue(stateBag,"Subnet.SubnetId",*resp).(string)
			s.subnetId = s.KingcloudRunConfig.SubnetId
		}
		return multistep.ActionContinue
	}
}

func (s *stepConfigKingcloudSubnet) Cleanup(stateBag multistep.StateBag) {
	if s.subnetId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Waiting Instance unbind Subnet "))
		_,waitErr := client.WaitSubnetClean(stateBag,s.subnetId)
		if waitErr != nil {
			ui.Error(fmt.Sprintf("Error waiting Subnet unbind %s", waitErr))
		}
		ui.Say(fmt.Sprintf("Deleting Subnet with Id %s ",s.subnetId))
		deleteSubnet := make(map[string]interface{})
		deleteSubnet["SubnetId"] = s.subnetId
		_,err := client.VpcClient.DeleteSubnet(&deleteSubnet)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete subnet %s", err))
		}
	}
}


