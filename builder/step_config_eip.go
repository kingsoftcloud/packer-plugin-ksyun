package ksyun

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"strconv"
)

type StepConfigKsyunPublicIp struct {
	CommonConfig *CommonConfig
	eipId        string
}

func (s *StepConfigKsyunPublicIp) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("ksyun_client").(*ClientWrapper)
	instanceId := stateBag.Get("InstanceId").(string)
	privateIp := stateBag.Get("PrivateIp").(string)
	chargeTypes := []string{"Daily", "TrafficMonthly", "DailyPaidByTransfer", "HourlyInstantSettlement"}
	checkChargeType := false
	if s.CommonConfig.AssociatePublicIpAddress {
		if s.CommonConfig.PublicIpBandWidth == 0 {
			// default bandwidth is 1 m
			s.CommonConfig.PublicIpBandWidth = 1
		} else if s.CommonConfig.PublicIpBandWidth > 100 {
			return Halt(stateBag, fmt.Errorf("public_ip max bandwidth must lower than 100"), "")
		}
		if s.CommonConfig.PublicIpChargeType == "" {
			// default PublicIpChargeType is Daily
			s.CommonConfig.PublicIpChargeType = "Daily"
			checkChargeType = true
		} else {
			for _, v := range chargeTypes {
				if s.CommonConfig.PublicIpChargeType == v {
					checkChargeType = true
				}
			}
		}
		if checkChargeType {
			ui.Say("Allocating eip...")
			//create eip
			createEip := make(map[string]interface{})
			createEip["BandWidth"] = strconv.Itoa(s.CommonConfig.PublicIpBandWidth)
			createEip["ChargeType"] = s.CommonConfig.PublicIpChargeType
			createEip["ProjectId"] = s.CommonConfig.ProjectId
			createResp, createErr := client.EipClient.AllocateAddress(&createEip)
			if createErr != nil {
				return Halt(stateBag, createErr, "Error creating new eip")
			}
			if createResp != nil {
				allocationId := GetSdkValue(stateBag, "AllocationId", *createResp).(string)
				publicIp := GetSdkValue(stateBag, "PublicIp", *createResp).(string)
				s.eipId = allocationId
				stateBag.Put("publicIp", publicIp)
				ui.Say("Associating eip to instance")
				//create_security_group_rule
				authorizeSecurityGroupEntry := make(map[string]interface{})
				authorizeSecurityGroupEntry["securityGroupId"] = s.CommonConfig.SecurityGroupId
				authorizeSecurityGroupEntry["CidrBlock"] = "0.0.0.0/0"
				authorizeSecurityGroupEntry["Direction"] = "in"
				authorizeSecurityGroupEntry["Protocol"] = "tcp"
				authorizeSecurityGroupEntry["PortRangeFrom"] = strconv.Itoa(22)
				authorizeSecurityGroupEntry["PortRangeTo"] = strconv.Itoa(22)
				_, errRule := client.VpcClient.AuthorizeSecurityGroupEntry(&authorizeSecurityGroupEntry)
				if errRule != nil {
					return Halt(stateBag, errRule, "Error creating  eip SecurityGroupRule")
				}
				//associate eip
				associateAddress := make(map[string]interface{})
				associateAddress["AllocationId"] = allocationId
				associateAddress["InstanceType"] = "Ipfwd"
				associateAddress["InstanceId"] = instanceId
				_, err := client.EipClient.AssociateAddress(&associateAddress)
				if err != nil {
					return Halt(stateBag, err, "Error associate eip to instance")
				}
			}

		} else {
			return Halt(stateBag, fmt.Errorf("public_ip_charge_type not match"), "")
		}
	} else {
		stateBag.Put("publicIp", privateIp)
	}
	return multistep.ActionContinue
}

func (s *StepConfigKsyunPublicIp) Cleanup(stateBag multistep.StateBag) {
	if s.eipId != "" {
		ui := stateBag.Get("ui").(packersdk.Ui)
		client := stateBag.Get("ksyun_client").(*ClientWrapper)
		ui.Say(fmt.Sprintf("Disassociate Eip with Id %s ", s.eipId))
		disassociateEip := make(map[string]interface{})
		disassociateEip["AllocationId"] = s.eipId
		_, disassociateErr := client.EipClient.DisassociateAddress(&disassociateEip)
		if disassociateErr != nil {
			ui.Error(fmt.Sprintf("Error disassociate Eip %s", disassociateErr))
		}
		ui.Say(fmt.Sprintf("Deleting Eip with Id %s ", s.eipId))
		deleteEip := make(map[string]interface{})
		deleteEip["AllocationId"] = s.eipId
		_, err := client.EipClient.ReleaseAddress(&deleteEip)
		if err != nil {
			ui.Error(fmt.Sprintf("Error delete Eip %s", err))
		}
	}
}
