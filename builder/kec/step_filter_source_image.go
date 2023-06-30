package kec

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

type stepFilterSourceImage struct {
	SourceImageId string
	KmiFilters    *ksyun.KmiFilterOptions
}

func (s *stepFilterSourceImage) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientKecWrapper)
	// instanceId := stateBag.Get("InstanceId").(string)
	// TODO: prepare information filtered
	ui.Say(fmt.Sprintf("%v", s))
	var params *map[string]interface{}
	if s.SourceImageId != "" {
		params = &map[string]interface{}{
			"ImageId": s.SourceImageId,
		}
	}
	image, err := s.KmiFilters.GetFilteredImage(params, client.KecClient)
	if err != nil {
		return ksyun.Halt(stateBag, err, "Error filtering source image")
	}

	ui.Say(fmt.Sprintf("Found an image, Id: %s, Platform %s", image.ImageId, image.Platform))

	stateBag.Put("source_image", image)
	return multistep.ActionContinue
}

func (s *stepFilterSourceImage) Cleanup(stateBag multistep.StateBag) {
}
