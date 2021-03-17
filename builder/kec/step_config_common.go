package kec

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type stepConfigKsyunCommon struct {
	KsyunRunConfig *KsyunRunConfig
}

func (s *stepConfigKsyunCommon) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	//set default project
	if s.KsyunRunConfig.ProjectId == "" {
		s.KsyunRunConfig.ProjectId = defaultProjectId
	}
	return multistep.ActionContinue
}

func (s *stepConfigKsyunCommon) Cleanup(stateBag multistep.StateBag) {

}
