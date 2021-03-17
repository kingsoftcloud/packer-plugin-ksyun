package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type stepCheckKingcloudSourceImage struct {
	SourceImageId string
}

func (s *stepCheckKingcloudSourceImage) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientWrapper)

	//query
	describeImages := make(map[string]interface{})
	describeImages["ImageId"] = s.SourceImageId
	_, err :=client.KecClient.DescribeImages(&describeImages)
	if err != nil {
		return Halt(stateBag,err,"Error querying kingcloud image")
	}

	ui.Message(fmt.Sprintf("Found image ID: %s", s.SourceImageId))
	return multistep.ActionContinue
}

func (s *stepCheckKingcloudSourceImage) Cleanup(bag multistep.StateBag) {
}


