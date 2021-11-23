//go:generate mapstructure-to-hcl2 -type Config
package epc

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
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

const BuilderId = "ksyun.epc"

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	ClientEpcConfig     `mapstructure:",squash"`
	KsyunEpcRunConfig   `mapstructure:",squash"`
	KsyunImageConfig    `mapstructure:",squash"`
	ctx                 interpolate.Context
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
	errs = packersdk.MultiErrorAppend(errs, b.config.ClientConfig.Prepare(&b.config.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, b.config.KsyunEpcRunConfig.Prepare(&b.config.ctx)...)
	errs = packersdk.MultiErrorAppend(errs, b.config.KsyunImageConfig.Prepare(&b.config.ctx)...)
	if errs != nil && len(errs.Errors) > 0 {
		return nil, nil, errs
	}

	packersdk.LogSecretFilter.Set(b.config.KsyunAccessKey, b.config.KsyunSecretKey)
	return nil, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packersdk.Ui, hook packersdk.Hook) (packersdk.Artifact, error) {
	stateBag := new(multistep.BasicStateBag)
	client := b.config.EpcClient(stateBag)
	stateBag.Put("config", &b.config)
	stateBag.Put("client", client)
	stateBag.Put("hook", hook)
	stateBag.Put("ui", ui)
	//special
	SSHTemporaryPublicKey := ""
	//step
	prepareSteps := []multistep.Step{
		&stepCheckKsyunSourceImage{
			SourceImageId: b.config.SourceImageId,
		},
		&ksyun.StepConfigKsyunKeyPair{
			CommonConfig:          &b.config.CommonConfig,
			SSHTemporaryPublicKey: &SSHTemporaryPublicKey,
		},
		&ksyun.StepConfigKsyunCommon{
			CommonConfig: &b.config.CommonConfig,
		},
		&ksyun.StepConfigKsyunVpc{
			CommonConfig: &b.config.CommonConfig,
		},
		&ksyun.StepConfigKsyunSubnet{
			CommonConfig: &b.config.CommonConfig,
			SubnetType:   "Physical",
			Index:        0,
			After:        b.config.PrepareExtensionSubnet(&b.config.CommonConfig),
		},
		&ksyun.StepConfigKsyunSecurityGroup{
			CommonConfig: &b.config.CommonConfig,
			Index:        0,
			After:        b.config.PrepareExtensionSecurityGroup(&b.config.CommonConfig),
		},
	}
	if b.config.NeedExtensionNetwork() {
		extraSteps := []multistep.Step{
			&ksyun.StepConfigKsyunSubnet{
				CommonConfig: &b.config.CommonConfig,
				SubnetType:   "Physical",
				Index:        1,
				After:        b.config.MergeExtensionSubnet(&b.config.CommonConfig),
			},
			&ksyun.StepConfigKsyunSecurityGroup{
				CommonConfig: &b.config.CommonConfig,
				Index:        1,
				After:        b.config.MergeExtensionSecurityGroup(&b.config.CommonConfig),
			},
		}
		prepareSteps = append(prepareSteps, extraSteps...)
	}
	runSteps := []multistep.Step{
		&stepCreateBareMetal{
			KsyunRunConfig: &b.config.KsyunEpcRunConfig,
		},
		&ksyun.StepConfigKsyunPublicIp{
			CommonConfig: &b.config.CommonConfig,
		},
		&communicator.StepConnect{
			Config:    &b.config.KsyunEpcRunConfig.Comm,
			Host:      ksyun.SSHHost(b.config.KsyunEpcRunConfig.Comm),
			SSHConfig: b.config.KsyunEpcRunConfig.Comm.SSHConfigFunc(),
		},
		&commonsteps.StepProvision{},
		&ksyun.StepCleanupKsyunTempKeys{
			Comm:                  &b.config.KsyunEpcRunConfig.Comm,
			SSHTemporaryPublicKey: &SSHTemporaryPublicKey,
		},
		&stepStopBareMetal{
			KsyunRunConfig: &b.config.KsyunEpcRunConfig,
		},
		&stepCreateKsyunImage{
			KsyunRunConfig:   &b.config.KsyunEpcRunConfig,
			KsyunImageConfig: &b.config.KsyunImageConfig,
		},
	}

	steps := append(prepareSteps, runSteps...)
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
		Client: b.config.client,
	}
	artifact.Artifact = &ksyun.Artifact{
		KsyunImageId:   stateBag.Get("TargetImageId").(string),
		BuilderIdValue: BuilderId,
	}
	return artifact, nil
}
