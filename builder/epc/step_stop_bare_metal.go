package epc

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type stepStopBareMetal struct {
	KsyunRunConfig *KsyunEpcRunConfig
}

func (s *stepStopBareMetal) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientEpcWrapper)
	instanceId := stateBag.Get("InstanceId").(string)

	ui.Say("Stopping Ksyun Bare Metal Instance ")
	stopInstance := make(map[string]interface{})
	stopInstance["HostId"] = instanceId
	_, errorStop := client.EpcClient.StopEpc(&stopInstance)
	if errorStop != nil {
		return ksyun.Halt(stateBag, errorStop, "Error stopping  Bare Metal instance")
	}
	ui.Say("Waiting Ksyun Kec Instance stopped ")
	_, err := client.WaitEpcInstanceStatus(stateBag, instanceId, s.KsyunRunConfig.ProjectId, "Stopped")
	if err != nil {
		return ksyun.Halt(stateBag, err, "Error waiting  Bare Metal instance status")
	}
	return multistep.ActionContinue
}

func (s *stepStopBareMetal) Cleanup(stateBag multistep.StateBag) {
}
