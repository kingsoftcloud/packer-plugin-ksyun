package kec

import (
	"context"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
)

type stepConfigKingcloudCommon struct {
	KingcloudRunConfig *KingcloudRunConfig
}

func (s *stepConfigKingcloudCommon) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	//set default project
	if s.KingcloudRunConfig.ProjectId == "" {
		s.KingcloudRunConfig.ProjectId = defaultProjectId
	}
	return multistep.ActionContinue
}

func (s *stepConfigKingcloudCommon) Cleanup(stateBag multistep.StateBag) {

}


