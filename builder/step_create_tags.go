package ksyun

import (
	"context"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

type StepCreateKsyunTags struct {
	Tags map[string]string

	// tag resource type Values: instance, epc-host, epc-image, image, eip
	ResourceType string
}

func (s *StepCreateKsyunTags) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("ksyun_client").(*ClientWrapper)
	conn := client.TagsClient
	targetImageId := stateBag.Get("TargetImageId").(string)

	if len(s.Tags) == 0 {
		return multistep.ActionContinue
	}

	ksyunTags := TagMap(s.Tags).KsyunTags()
	ksyunTags.Report(ui)

	if !StringInSlice(s.ResourceType, ResourceTypeList, false) {
		s.ResourceType = ResourceTypeImage
	}
	processTagsParams, err := ksyunTags.GetTagsParams(s.ResourceType, targetImageId)
	if err != nil {
		return Halt(stateBag, err, "Error processing tag")
	}
	_, err = conn.ReplaceResourcesTags(&processTagsParams)
	if err != nil {
		return Halt(stateBag, err, "Error bounding tag")
	}

	return multistep.ActionContinue
}

func (s *StepCreateKsyunTags) Cleanup(stateBag multistep.StateBag) {
	// No cleanup...
}
