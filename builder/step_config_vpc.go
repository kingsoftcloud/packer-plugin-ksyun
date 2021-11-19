package ksyun

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepConfigKsyunVpc struct {
	CommonConfig *CommonConfig
	vpcId        string
}

func (s *StepConfigKsyunVpc) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("ksyun_client").(*ClientWrapper)

	if s.CommonConfig.VpcId != "" {
		//Check_vpc
		queryVpc := make(map[string]interface{})
		queryVpc["VpcId.1"] = s.CommonConfig.VpcId
		resp, err := client.VpcClient.DescribeVpcs(&queryVpc)
		if err != nil {
			return Halt(stateBag, err, fmt.Sprintf("Error query Vpc with id %s", s.CommonConfig.VpcId))
		}
		if resp != nil {
			vpcId := GetSdkValue(stateBag, "VpcSet.0.VpcId", *resp)
			if vpcId == nil {
				return Halt(stateBag, err, fmt.Sprintf("Vpc id %s not found", s.CommonConfig.VpcId))
			}
		}
		ui.Say(fmt.Sprintf("Using existing Vpc id is %s", s.CommonConfig.VpcId))
		return multistep.ActionContinue
	} else {
		//create_vpc
		if s.CommonConfig.VpcName == "" {
			s.CommonConfig.VpcName = defaultVpcName
		}
		if s.CommonConfig.VpcCidrBlock == "" {
			s.CommonConfig.VpcCidrBlock = defaultVpcCidr
		}
		ui.Say(fmt.Sprintf("Creating new Vpc with name %s cidr %s", s.CommonConfig.VpcName,
			s.CommonConfig.VpcCidrBlock))
		createVpc := make(map[string]interface{})
		createVpc["VpcName"] = s.CommonConfig.VpcName
		createVpc["CidrBlock"] = s.CommonConfig.VpcCidrBlock
		resp, err := client.VpcClient.CreateVpc(&createVpc)
		if err != nil {
			return Halt(stateBag, err, "Error creating new Vpc")
		}
		if resp != nil {
			s.CommonConfig.VpcId = GetSdkValue(stateBag, "Vpc.VpcId", *resp).(string)
			s.vpcId = s.CommonConfig.VpcId
		}
		return multistep.ActionContinue
	}
}

func (s *StepConfigKsyunVpc) Cleanup(stateBag multistep.StateBag) {
	if s.vpcId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("ksyun_client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Vpc with Id %s ", s.vpcId))
		deleteVpc := make(map[string]interface{})
		deleteVpc["VpcId"] = s.vpcId
		_, err := client.VpcClient.DeleteVpc(&deleteVpc)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete vpc %s", err))
		}
	}
}
