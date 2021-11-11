package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type stepConfigKsyunSubnet struct {
	KsyunRunConfig *KsyunKecRunConfig
	subnetId       string
}

func (s *stepConfigKsyunSubnet) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)

	if s.KsyunRunConfig.SubnetId != "" {
		//Check_Subnet
		querySubnet := make(map[string]interface{})
		querySubnet["SubnetId.1"] = s.KsyunRunConfig.SubnetId
		resp, err := client.VpcClient.DescribeSubnets(&querySubnet)
		if err != nil {
			return ksyun.Halt(stateBag, err, fmt.Sprintf("Error query Subnet with id %s", s.KsyunRunConfig.SubnetId))
		}
		if resp != nil {
			subnetId := ksyun.GetSdkValue(stateBag, "SubnetSet.0.SubnetId", *resp)
			subnetType := ksyun.GetSdkValue(stateBag, "SubnetSet.0.SubnetType", *resp)
			subnetName := ksyun.GetSdkValue(stateBag, "SubnetSet.0.SubnetName", *resp)
			vpcId := ksyun.GetSdkValue(stateBag, "SubnetSet.0.VpcId", *resp)
			if subnetId == nil {
				return ksyun.Halt(stateBag, err, fmt.Sprintf("Subnet id %s not found", s.KsyunRunConfig.SubnetId))
			}

			if vpcId != s.KsyunRunConfig.VpcId {
				return ksyun.Halt(stateBag, fmt.Errorf(fmt.Sprintf("Subnet id %s vpc not match",
					s.KsyunRunConfig.SubnetId)), "")
			}

			if subnetType != EnableSubnetType {
				return ksyun.Halt(stateBag,
					fmt.Errorf(fmt.Sprintf("Subnet id %s Type is Not %s", EnableSubnetType, s.KsyunRunConfig.SubnetId)), "")
			}
			ui.Say(fmt.Sprintf("Using existing Subnet id is %s name is %s", s.KsyunRunConfig.SubnetId,
				subnetName))
		}
		return multistep.ActionContinue
	} else {
		//create_subnet
		if s.KsyunRunConfig.SubnetName == "" {
			s.KsyunRunConfig.SubnetName = defaultSubnetName
		}
		if s.KsyunRunConfig.SubnetCidrBlock == "" {
			s.KsyunRunConfig.SubnetCidrBlock = defaultSubnetCidr
		}
		startIp, minIp, maxIp := ksyun.GetCidrIpRange(s.KsyunRunConfig.SubnetCidrBlock)
		ui.Say(fmt.Sprintf("Creating new Subnet with name  %s cidr %s vpcId %s",
			s.KsyunRunConfig.SubnetName, s.KsyunRunConfig.SubnetCidrBlock, s.KsyunRunConfig.VpcId))
		createSubnet := make(map[string]interface{})
		createSubnet["VpcId"] = s.KsyunRunConfig.VpcId
		createSubnet["SubnetName"] = s.KsyunRunConfig.SubnetName
		createSubnet["SubnetType"] = EnableSubnetType
		createSubnet["CidrBlock"] = s.KsyunRunConfig.SubnetCidrBlock
		createSubnet["GatewayIp"] = startIp
		createSubnet["DhcpIpFrom"] = minIp
		createSubnet["DhcpIpTo"] = maxIp
		if s.KsyunRunConfig.AvailabilityZone != "" {
			createSubnet["AvailabilityZone"] = s.KsyunRunConfig.AvailabilityZone
		}
		resp, err := client.VpcClient.CreateSubnet(&createSubnet)
		if err != nil {
			return ksyun.Halt(stateBag, err, "Error creating new Subnet")
		}
		if resp != nil {
			s.KsyunRunConfig.SubnetId = ksyun.GetSdkValue(stateBag, "Subnet.SubnetId", *resp).(string)
			s.subnetId = s.KsyunRunConfig.SubnetId
		}
		return multistep.ActionContinue
	}
}

func (s *stepConfigKsyunSubnet) Cleanup(stateBag multistep.StateBag) {
	if s.subnetId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Waiting Instance unbind Subnet "))
		_, waitErr := client.WaitSubnetClean(stateBag, s.subnetId)
		if waitErr != nil {
			ui.Error(fmt.Sprintf("Error waiting Subnet unbind %s", waitErr))
		}
		ui.Say(fmt.Sprintf("Deleting Subnet with Id %s ", s.subnetId))
		deleteSubnet := make(map[string]interface{})
		deleteSubnet["SubnetId"] = s.subnetId
		_, err := client.VpcClient.DeleteSubnet(&deleteSubnet)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete subnet %s", err))
		}
	}
}
