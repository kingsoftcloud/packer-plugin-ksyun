package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepConfigKsyunSecurityGroup struct {
	KsyunRunConfig  *KsyunRunConfig
	SecurityGroupId string
}

func (s *stepConfigKsyunSecurityGroup) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {

	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)

	if s.KsyunRunConfig.SecurityGroupId != "" {
		//Check_security_group
		querySecurityGroup := make(map[string]interface{})
		querySecurityGroup["SecurityGroupId.1"] = s.KsyunRunConfig.SecurityGroupId
		resp, err := client.VpcClient.DescribeSubnets(&querySecurityGroup)
		if err != nil {
			return Halt(stateBag, err, fmt.Sprintf("Error query SecurityGroup with id %s",
				s.KsyunRunConfig.SecurityGroupId))
		}
		if resp != nil {
			securityGroupId := getSdkValue(stateBag, "SecurityGroupSet.0.SecurityGroupId", *resp)
			securityGroupName := getSdkValue(stateBag, "SecurityGroupSet.0.SecurityGroupName", *resp)
			vpcId := getSdkValue(stateBag, "SecurityGroupSet.0.VpcId", *resp)
			if securityGroupId == nil {
				return Halt(stateBag, fmt.Errorf(fmt.Sprintf("SecurityGroup id %s not found",
					s.KsyunRunConfig.SecurityGroupId)), "")
			}

			if vpcId != s.KsyunRunConfig.VpcId {
				return Halt(stateBag, fmt.Errorf(fmt.Sprintf("SecurityGroup id %s vpc not match",
					s.KsyunRunConfig.SecurityGroupId)), "")
			}

			ui.Say(fmt.Sprintf("Using existing SecurityGroup id is %s name is %s ",
				s.KsyunRunConfig.SecurityGroupId, securityGroupName))
		}
		return multistep.ActionContinue
	} else {
		//create_security_group
		if s.KsyunRunConfig.SecurityGroupName == "" {
			s.KsyunRunConfig.SecurityGroupName = defaultSecurityGroupName
		}
		ui.Say(fmt.Sprintf("Creating new SecurityGroup with name  %s vpcId %s",
			s.KsyunRunConfig.SecurityGroupName, s.KsyunRunConfig.VpcId))
		createSecurityGroup := make(map[string]interface{})
		createSecurityGroup["VpcId"] = s.KsyunRunConfig.VpcId
		createSecurityGroup["SecurityGroupName"] = s.KsyunRunConfig.SecurityGroupName
		resp, err := client.VpcClient.CreateSecurityGroup(&createSecurityGroup)
		if err != nil {
			return Halt(stateBag, err, "Error creating new SecurityGroup")
		}
		if resp != nil {
			s.KsyunRunConfig.SecurityGroupId = getSdkValue(stateBag, "SecurityGroup.SecurityGroupId", *resp).(string)
			s.SecurityGroupId = s.KsyunRunConfig.SecurityGroupId
			//create_security_group_rule
			//authorizeSecurityGroupEntry := make(map[string]interface{})
			//authorizeSecurityGroupEntry["SecurityGroupId"] = s.KsyunRunConfig.SecurityGroupId
			//authorizeSecurityGroupEntry["CidrBlock"] = "0.0.0.0/0"
			//authorizeSecurityGroupEntry["Direction"] = "out"
			//authorizeSecurityGroupEntry["Protocol"] = "ip"
			//_,err1 := client.VpcClient.AuthorizeSecurityGroupEntry(&authorizeSecurityGroupEntry)
			//if err1 != nil {
			//	return Halt(stateBag, err1, "Error creating new SecurityGroupRule")
			//}

		}
		return multistep.ActionContinue
	}
}

func (s *stepConfigKsyunSecurityGroup) Cleanup(stateBag multistep.StateBag) {
	if s.SecurityGroupId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Waiting Instance unbind SecurityGroup "))
		_, waitErr := client.WaitSecurityGroupClean(stateBag, s.SecurityGroupId)
		if waitErr != nil {
			ui.Error(fmt.Sprintf("Error waiting SecurityGroup unbind %s", waitErr))
		}
		ui.Say(fmt.Sprintf("Deleting SecurityGroup with Id %s ", s.SecurityGroupId))
		deleteSecurityGroup := make(map[string]interface{})
		deleteSecurityGroup["SecurityGroupId"] = s.SecurityGroupId
		_, err := client.VpcClient.DeleteSecurityGroup(&deleteSecurityGroup)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete SecurityGroup %s", err))
		}
	}
}
