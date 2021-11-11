package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type stepConfigKsyunVpc struct {
	KsyunRunConfig *KsyunKecRunConfig
	vpcId          string
}

func (s *stepConfigKsyunVpc) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)

	if s.KsyunRunConfig.VpcId != "" {
		//Check_vpc
		queryVpc := make(map[string]interface{})
		queryVpc["VpcId.1"] = s.KsyunRunConfig.VpcId
		resp, err := client.VpcClient.DescribeVpcs(&queryVpc)
		if err != nil {
			return ksyun.Halt(stateBag, err, fmt.Sprintf("Error query Vpc with id %s", s.KsyunRunConfig.VpcId))
		}
		if resp != nil {
			vpcId := ksyun.GetSdkValue(stateBag, "VpcSet.0.VpcId", *resp)
			if vpcId == nil {
				return ksyun.Halt(stateBag, err, fmt.Sprintf("Vpc id %s not found", s.KsyunRunConfig.VpcId))
			}
		}
		ui.Say(fmt.Sprintf("Using existing Vpc id is %s", s.KsyunRunConfig.VpcId))
		return multistep.ActionContinue
	} else {
		//create_vpc
		if s.KsyunRunConfig.VpcName == "" {
			s.KsyunRunConfig.VpcName = defaultVpcName
		}
		if s.KsyunRunConfig.VpcCidrBlock == "" {
			s.KsyunRunConfig.VpcCidrBlock = defaultVpcCidr
		}
		ui.Say(fmt.Sprintf("Creating new Vpc with name %s cidr %s", s.KsyunRunConfig.VpcName,
			s.KsyunRunConfig.VpcCidrBlock))
		createVpc := make(map[string]interface{})
		createVpc["VpcName"] = s.KsyunRunConfig.VpcName
		createVpc["CidrBlock"] = s.KsyunRunConfig.VpcCidrBlock
		resp, err := client.VpcClient.CreateVpc(&createVpc)
		if err != nil {
			return ksyun.Halt(stateBag, err, "Error creating new Vpc")
		}
		if resp != nil {
			s.KsyunRunConfig.VpcId = ksyun.GetSdkValue(stateBag, "Vpc.VpcId", *resp).(string)
			s.vpcId = s.KsyunRunConfig.VpcId
		}
		return multistep.ActionContinue
	}
}

func (s *stepConfigKsyunVpc) Cleanup(stateBag multistep.StateBag) {
	if s.vpcId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Vpc with Id %s ", s.vpcId))
		deleteVpc := make(map[string]interface{})
		deleteVpc["VpcId"] = s.vpcId
		_, err := client.VpcClient.DeleteVpc(&deleteVpc)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete vpc %s", err))
		}
	}
}
