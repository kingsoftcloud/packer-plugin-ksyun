package kec

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepStopKsyunKec struct {
	KsyunRunConfig *KsyunRunConfig
}

func (s *stepStopKsyunKec) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)
	instanceId := stateBag.Get("InstanceId").(string)

	ui.Say("Stopping Ksyun Kec Instance ")
	stopInstance := make(map[string]interface{})
	stopInstance["InstanceId.1"] = instanceId
	_, errorStop := client.KecClient.StopInstances(&stopInstance)
	if errorStop != nil {
		return Halt(stateBag, errorStop, "Error stopping  kec instance")
	}
	ui.Say("Waiting Ksyun Kec Instance stopped ")
	_, err := client.WaitKecInstanceStatus(stateBag, instanceId, s.KsyunRunConfig.ProjectId, "stopped")
	if err != nil {
		return Halt(stateBag, err, "Error waiting  kec instance status")
	}
	return multistep.ActionContinue
}

func (s *stepStopKsyunKec) Cleanup(stateBag multistep.StateBag) {
}
