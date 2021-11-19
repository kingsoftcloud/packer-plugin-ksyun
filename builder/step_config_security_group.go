package ksyun

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepConfigKsyunSecurityGroup struct {
	CommonConfig    *CommonConfig
	securityGroupId string
}

func (s *StepConfigKsyunSecurityGroup) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {

	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("ksyun_client").(*ClientWrapper)

	if s.CommonConfig.SecurityGroupId != "" {
		//Check_security_group
		querySecurityGroup := make(map[string]interface{})
		querySecurityGroup["securityGroupId.1"] = s.CommonConfig.SecurityGroupId
		resp, err := client.VpcClient.DescribeSubnets(&querySecurityGroup)
		if err != nil {
			return Halt(stateBag, err, fmt.Sprintf("Error query SecurityGroup with id %s",
				s.CommonConfig.SecurityGroupId))
		}
		if resp != nil {
			securityGroupId := GetSdkValue(stateBag, "SecurityGroupSet.0.securityGroupId", *resp)
			securityGroupName := GetSdkValue(stateBag, "SecurityGroupSet.0.SecurityGroupName", *resp)
			vpcId := GetSdkValue(stateBag, "SecurityGroupSet.0.VpcId", *resp)
			if securityGroupId == nil {
				return Halt(stateBag, fmt.Errorf(fmt.Sprintf("SecurityGroup id %s not found",
					s.CommonConfig.SecurityGroupId)), "")
			}

			if vpcId != s.CommonConfig.VpcId {
				return Halt(stateBag, fmt.Errorf(fmt.Sprintf("SecurityGroup id %s vpc not match",
					s.CommonConfig.SecurityGroupId)), "")
			}

			ui.Say(fmt.Sprintf("Using existing SecurityGroup id is %s name is %s ",
				s.CommonConfig.SecurityGroupId, securityGroupName))
		}
		return multistep.ActionContinue
	} else {
		//create_security_group
		if s.CommonConfig.SecurityGroupName == "" {
			s.CommonConfig.SecurityGroupName = defaultSecurityGroupName
		}
		ui.Say(fmt.Sprintf("Creating new SecurityGroup with name  %s vpcId %s",
			s.CommonConfig.SecurityGroupName, s.CommonConfig.VpcId))
		createSecurityGroup := make(map[string]interface{})
		createSecurityGroup["VpcId"] = s.CommonConfig.VpcId
		createSecurityGroup["SecurityGroupName"] = s.CommonConfig.SecurityGroupName
		resp, err := client.VpcClient.CreateSecurityGroup(&createSecurityGroup)
		if err != nil {
			return Halt(stateBag, err, "Error creating new SecurityGroup")
		}
		if resp != nil {
			s.CommonConfig.SecurityGroupId = GetSdkValue(stateBag, "SecurityGroup.SecurityGroupId", *resp).(string)
			s.securityGroupId = s.CommonConfig.SecurityGroupId
			//create_security_group_rule
			//authorizeSecurityGroupEntry := make(map[string]interface{})
			//authorizeSecurityGroupEntry["securityGroupId"] = s.KsyunKecRunConfig.securityGroupId
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

func (s *StepConfigKsyunSecurityGroup) Cleanup(stateBag multistep.StateBag) {
	if s.securityGroupId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("ksyun_client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Waiting Instance unbind SecurityGroup "))
		_, waitErr := client.WaitSecurityGroupClean(stateBag, s.securityGroupId)
		if waitErr != nil {
			ui.Error(fmt.Sprintf("Error waiting SecurityGroup unbind %s", waitErr))
		}
		ui.Say(fmt.Sprintf("Deleting SecurityGroup with Id %s ", s.securityGroupId))
		deleteSecurityGroup := make(map[string]interface{})
		deleteSecurityGroup["SecurityGroupId"] = s.securityGroupId
		_, err := client.VpcClient.DeleteSecurityGroup(&deleteSecurityGroup)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete SecurityGroup %s", err))
		}
	}
}
