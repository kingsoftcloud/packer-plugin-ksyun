package epc

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type stepCreateKsyunImage struct {
	KsyunRunConfig   *KsyunEpcRunConfig
	KsyunImageConfig *KsyunImageConfig
}

func (s *stepCreateKsyunImage) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientEpcWrapper)
	instanceId := stateBag.Get("InstanceId").(string)
	ui.Say("Creating Ksyun Bare Metal Image ")
	createImage := make(map[string]interface{})
	createImage["HostId"] = instanceId
	createImage["ImageName"] = s.KsyunImageConfig.KsyunImageName
	resp, errorCreate := client.EpcClient.CreateImage(&createImage)
	if errorCreate != nil {
		return ksyun.Halt(stateBag, errorCreate, "Error creating  Bare Metal image")
	}
	if resp != nil {
		ui.Say("Waiting Ksyun Bare Metal Image active")
		imageId := ksyun.GetSdkValue(stateBag, "Image.ImageId", *resp).(string)
		_, err := client.WaitEpcInstanceStatus(stateBag, instanceId, s.KsyunRunConfig.ProjectId, "Running")
		if err != nil {
			return ksyun.Halt(stateBag, err, "Error waiting  Bare Metal image active")
		}
		stateBag.Put("TargetImageId", imageId)
	}
	return multistep.ActionContinue
}

func (s *stepCreateKsyunImage) Cleanup(stateBag multistep.StateBag) {

}
