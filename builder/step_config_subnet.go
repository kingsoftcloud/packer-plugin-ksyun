package ksyun

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"strconv"
	"strings"
)

type StepConfigKsyunSubnet struct {
	CommonConfig *CommonConfig
	subnetId     string
	SubnetType   string
	Index        int
	After        AfterStepRun
}

func (s *StepConfigKsyunSubnet) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	defer func() {
		if s.After != nil {
			s.After()
		}
	}()
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("ksyun_client").(*ClientWrapper)
	if s.CommonConfig.SubnetId != "" {
		//Check_Subnet
		querySubnet := make(map[string]interface{})
		querySubnet["SubnetId.1"] = s.CommonConfig.SubnetId
		resp, err := client.VpcClient.DescribeSubnets(&querySubnet)
		if err != nil {
			return Halt(stateBag, err, fmt.Sprintf("Error query Subnet with id %s", s.CommonConfig.SubnetId))
		}
		if resp != nil {
			subnetId := GetSdkValue(stateBag, "SubnetSet.0.SubnetId", *resp)
			subnetType := GetSdkValue(stateBag, "SubnetSet.0.SubnetType", *resp)
			subnetName := GetSdkValue(stateBag, "SubnetSet.0.SubnetName", *resp)
			vpcId := GetSdkValue(stateBag, "SubnetSet.0.VpcId", *resp)
			if subnetId == nil {
				return Halt(stateBag, err, fmt.Sprintf("Subnet id %s not found", s.CommonConfig.SubnetId))
			}

			if vpcId != s.CommonConfig.VpcId {
				return Halt(stateBag, fmt.Errorf(fmt.Sprintf("Subnet id %s vpc not match",
					s.CommonConfig.SubnetId)), "")
			}

			if subnetType != s.SubnetType {
				return Halt(stateBag,
					fmt.Errorf(fmt.Sprintf("Subnet id %s Type is Not %s", s.SubnetType, s.CommonConfig.SubnetId)), "")
			}
			s.CommonConfig.AvailabilityZone = GetSdkValue(stateBag, "SubnetSet.0.AvailabilityZoneName", *resp).(string)
			ui.Say(fmt.Sprintf("Using existing Subnet id is %s name is %s", s.CommonConfig.SubnetId,
				subnetName))
		}
		return multistep.ActionContinue
	} else {
		//create_subnet
		if s.CommonConfig.SubnetName == "" {
			s.CommonConfig.SubnetName = strings.Replace(defaultSubnetName, "index", strconv.Itoa(s.Index+1), -1)
		}
		if s.CommonConfig.SubnetCidrBlock == "" {
			s.CommonConfig.SubnetCidrBlock = strings.Replace(defaultSubnetCidr, "index", strconv.Itoa(s.Index+1), -1)
		}
		startIp, minIp, maxIp := GetCidrIpRange(s.CommonConfig.SubnetCidrBlock)
		ui.Say(fmt.Sprintf("Creating new Subnet with name  %s cidr %s vpcId %s",
			s.CommonConfig.SubnetName, s.CommonConfig.SubnetCidrBlock, s.CommonConfig.VpcId))
		createSubnet := make(map[string]interface{})
		createSubnet["VpcId"] = s.CommonConfig.VpcId
		createSubnet["SubnetName"] = s.CommonConfig.SubnetName
		createSubnet["SubnetType"] = s.SubnetType
		createSubnet["CidrBlock"] = s.CommonConfig.SubnetCidrBlock
		if s.CommonConfig.DNS1 != "" {
			createSubnet["DNS1"] = s.CommonConfig.DNS1
		}
		if s.CommonConfig.DNS2 != "" {
			createSubnet["DNS2"] = s.CommonConfig.DNS2
		}
		createSubnet["GatewayIp"] = startIp
		createSubnet["DhcpIpFrom"] = minIp
		createSubnet["DhcpIpTo"] = maxIp
		if s.CommonConfig.AvailabilityZone != "" {
			createSubnet["AvailabilityZone"] = s.CommonConfig.AvailabilityZone
		}
		resp, err := client.VpcClient.CreateSubnet(&createSubnet)

		if err != nil {
			return Halt(stateBag, err, "Error creating new Subnet")
		}
		if resp != nil {
			s.CommonConfig.SubnetId = GetSdkValue(stateBag, "Subnet.SubnetId", *resp).(string)
			s.CommonConfig.AvailabilityZone = GetSdkValue(stateBag, "Subnet.AvailabilityZoneName", *resp).(string)
			s.subnetId = s.CommonConfig.SubnetId
		}
		return multistep.ActionContinue
	}
}

func (s *StepConfigKsyunSubnet) Cleanup(stateBag multistep.StateBag) {
	if s.subnetId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("ksyun_client").(*ClientWrapper)
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
