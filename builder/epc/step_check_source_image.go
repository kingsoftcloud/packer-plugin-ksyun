package epc

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type stepCheckKsyunSourceImage struct {
	SourceImageId string
}

func (s *stepCheckKsyunSourceImage) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientEpcWrapper)

	//query
	describeImages := make(map[string]interface{})
	describeImages["ImageId.1"] = s.SourceImageId
	resp, err := client.EpcClient.DescribeImages(&describeImages)
	if err != nil {
		return ksyun.Halt(stateBag, err, "Error querying ksyun image")
	}
	if len(ksyun.GetSdkValue(stateBag, "ImageSet", *resp).([]interface{})) == 0 {
		return ksyun.Halt(stateBag, err, "ksyun image not found")
	}

	ui.Message(fmt.Sprintf("Found image ID: %s", s.SourceImageId))
	return multistep.ActionContinue
}

func (s *stepCheckKsyunSourceImage) Cleanup(bag multistep.StateBag) {
}
