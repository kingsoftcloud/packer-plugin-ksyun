//go:generate mapstructure-to-hcl2 -type Config,KingcloudDiskDevice,KingcloudEbsDataDisk

// The kingcloud  contains a packersdk.Builder implementation that
// builds ecs images for kingcloud.
package kec

import (
	"context"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

// The unique ID for this builder
const BuilderId = "kingcloud.kec"

type Config struct {
	common.PackerConfig   `mapstructure:",squash"`
	KingcloudAccessConfig `mapstructure:",squash"`
	KingcloudImageConfig  `mapstructure:",squash"`
	KingcloudRunConfig    `mapstructure:",squash"`

	ctx interpolate.Context
}

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	err := config.Decode(&b.config, &config.DecodeOpts{
		PluginType:         BuilderId,
		Interpolate:        true,
		InterpolateContext: &b.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{
				"run_command",
			},
		},
	}, raws...)
	b.config.ctx.EnableEnv = true
	if err != nil {
		return nil, nil, err
	}

	// Accumulate any errors
	var errs *packersdk.MultiError
	errs = packersdk.MultiErrorAppend(errs, b.config.KingcloudAccessConfig.Prepare(&b.config.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, b.config.KingcloudImageConfig.Prepare(&b.config.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, b.config.KingcloudRunConfig.Prepare(&b.config.ctx)...)

	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	packersdk.LogSecretFilter.Set(b.config.KingcloudAccessKey, b.config.KingcloudSecretKey)
	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packersdk.Ui, hook packersdk.Hook) (packersdk.Artifact, error) {
	client := b.config.Client()
	stateBag := new(multistep.BasicStateBag)
	stateBag.Put("config", &b.config)
	stateBag.Put("client", client)
	stateBag.Put("hook", hook)
	stateBag.Put("ui", ui)

	//step
	var steps []multistep.Step
	steps = []multistep.Step{
		&stepConfigKingcloudCommon{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
		},
		&stepCheckKingcloudSourceImage{
			SourceImageId: b.config.SourceImageId,
		},
		&stepConfigKingcloudKeyPair{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
			Comm: &b.config.KingcloudRunConfig.Comm,
		},
		&stepConfigKingcloudVpc{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
		},
		&stepConfigKingcloudSubnet{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
		},
		&stepConfigKingcloudSecurityGroup{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
		},
		&stepCreateKingcloudKec{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
		},
		&stepConfigKingcloudPublicIp{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
		},
		&communicator.StepConnect{
			Config: &b.config.KingcloudRunConfig.Comm,
			Host: SSHHost(b.config.KingcloudRunConfig.Comm),
			SSHConfig: b.config.KingcloudRunConfig.Comm.SSHConfigFunc(),
		},
		&commonsteps.StepProvision{},
		&stepStopKingcloudKec{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
		},
		&stepCreateKingcloudImage{
			KingcloudRunConfig: &b.config.KingcloudRunConfig,
			KingcloudImageConfig: &b.config.KingcloudImageConfig,
		},
	}



	// Run!
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, stateBag)

	// Build the artifact and return it
	artifact := &Artifact{
		KingcloudImageId: stateBag.Get("TargetImageId").(string),
		BuilderIdValue: BuilderId,
		Client:         b.config.client,
	}

	return artifact, nil
}