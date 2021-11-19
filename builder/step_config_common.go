package ksyun

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type StepConfigKsyunCommon struct {
	CommonConfig *CommonConfig
}

func (s *StepConfigKsyunCommon) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	//set default project
	if s.CommonConfig.ProjectId == "" {
		s.CommonConfig.ProjectId = defaultProjectId
	}
	return multistep.ActionContinue
}

func (s *StepConfigKsyunCommon) Cleanup(stateBag multistep.StateBag) {

}
