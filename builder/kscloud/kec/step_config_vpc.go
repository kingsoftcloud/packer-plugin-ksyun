package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepConfigKingcloudVpc struct {
	KingcloudRunConfig *KingcloudRunConfig
	vpcId string
}

func (s *stepConfigKingcloudVpc) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)

	if s.KingcloudRunConfig.VpcId != ""{
		//Check_vpc
		queryVpc := make(map[string]interface{})
		queryVpc["VpcId.1"] = s.KingcloudRunConfig.VpcId
		resp,err := client.VpcClient.DescribeVpcs(&queryVpc)
		if err != nil {
			return Halt(stateBag, err, fmt.Sprintf("Error query Vpc with id %s",s.KingcloudRunConfig.VpcId))
		}
		if resp !=nil{
			vpcId := getSdkValue(stateBag,"VpcSet.0.VpcId",*resp)
			if vpcId == nil {
				return Halt(stateBag, err, fmt.Sprintf("Vpc id %s not found",s.KingcloudRunConfig.VpcId))
			}
		}
		ui.Say(fmt.Sprintf("Using existing Vpc id is %s", s.KingcloudRunConfig.VpcId))
		return multistep.ActionContinue
	}else{
		//create_vpc
		if s.KingcloudRunConfig.VpcName == ""{
			s.KingcloudRunConfig.VpcName = defaultVpcName
		}
		if s.KingcloudRunConfig.VpcCidrBlock == ""{
			s.KingcloudRunConfig.VpcCidrBlock = defaultVpcCidr
		}
		ui.Say(fmt.Sprintf("Creating new Vpc with name %s cidr %s",s.KingcloudRunConfig.VpcName,
			s.KingcloudRunConfig.VpcCidrBlock))
		createVpc := make(map[string]interface{})
		createVpc["VpcName"] = s.KingcloudRunConfig.VpcName
		createVpc["CidrBlock"] = s.KingcloudRunConfig.VpcCidrBlock
		resp,err := client.VpcClient.CreateVpc(&createVpc)
		if err != nil {
			return Halt(stateBag, err, "Error creating new Vpc")
		}
		if resp !=nil {
			s.KingcloudRunConfig.VpcId = getSdkValue(stateBag,"Vpc.VpcId",*resp).(string)
			s.vpcId = s.KingcloudRunConfig.VpcId
		}
		return multistep.ActionContinue
	}
}

func (s *stepConfigKingcloudVpc) Cleanup(stateBag multistep.StateBag) {
	if s.vpcId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Deleting Vpc with Id %s ",s.vpcId))
		deleteVpc := make(map[string]interface{})
		deleteVpc["VpcId"] = s.vpcId
		_,err := client.VpcClient.DeleteVpc(&deleteVpc)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete vpc %s", err))
		}
	}
}



