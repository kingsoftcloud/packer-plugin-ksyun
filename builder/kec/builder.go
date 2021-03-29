//go:generate mapstructure-to-hcl2 -type Config,KsyunDiskDevice,KsyunEbsDataDisk

// The ksyun  contains a packersdk.Builder implementation that
// builds ecs images for ksyun.
package kec

import (
	"context"
	"fmt"
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
const BuilderId = "ksyun.kec"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	KsyunAccessConfig   `mapstructure:",squash"`
	KsyunImageConfig    `mapstructure:",squash"`
	KsyunRunConfig      `mapstructure:",squash"`

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
	errs = packersdk.MultiErrorAppend(errs, b.config.KsyunAccessConfig.Prepare(&b.config.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, b.config.KsyunImageConfig.Prepare(&b.config.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, b.config.KsyunRunConfig.Prepare(&b.config.ctx)...)

	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	packersdk.LogSecretFilter.Set(b.config.KsyunAccessKey, b.config.KsyunSecretKey)
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
		&stepConfigKsyunCommon{
			KsyunRunConfig: &b.config.KsyunRunConfig,
		},
		&stepCheckKsyunSourceImage{
			SourceImageId: b.config.SourceImageId,
		},
		&stepConfigKsyunKeyPair{
			KsyunRunConfig: &b.config.KsyunRunConfig,
			Comm:           &b.config.KsyunRunConfig.Comm,
		},
		&stepConfigKsyunVpc{
			KsyunRunConfig: &b.config.KsyunRunConfig,
		},
		&stepConfigKsyunSubnet{
			KsyunRunConfig: &b.config.KsyunRunConfig,
		},
		&stepConfigKsyunSecurityGroup{
			KsyunRunConfig: &b.config.KsyunRunConfig,
		},
		&stepCreateKsyunKec{
			KsyunRunConfig: &b.config.KsyunRunConfig,
		},
		&stepConfigKsyunPublicIp{
			KsyunRunConfig: &b.config.KsyunRunConfig,
		},
		&communicator.StepConnect{
			Config:    &b.config.KsyunRunConfig.Comm,
			Host:      SSHHost(b.config.KsyunRunConfig.Comm),
			SSHConfig: b.config.KsyunRunConfig.Comm.SSHConfigFunc(),
		},
		&commonsteps.StepProvision{},
		&commonsteps.StepCleanupTempKeys{
			Comm: &b.config.KsyunRunConfig.Comm,
		},
		&stepStopKsyunKec{
			KsyunRunConfig: &b.config.KsyunRunConfig,
		},
		&stepCreateKsyunImage{
			KsyunRunConfig:   &b.config.KsyunRunConfig,
			KsyunImageConfig: &b.config.KsyunImageConfig,
		},
	}

	// Run!
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, stateBag)

	// If there was an error, return that
	if err, ok := stateBag.GetOk("error"); ok {
		ui.Say(fmt.Sprintf("find some error %v ", err))
		return nil, err.(error)
	}

	// Build the artifact and return it
	artifact := &Artifact{
		KsyunImageId:   stateBag.Get("TargetImageId").(string),
		BuilderIdValue: BuilderId,
		Client:         b.config.client,
	}

	return artifact, nil
}
