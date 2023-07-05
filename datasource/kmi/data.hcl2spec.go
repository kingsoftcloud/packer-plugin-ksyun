// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package kmi

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	PackerBuildName     *string           `mapstructure:"packer_build_name" cty:"packer_build_name" hcl:"packer_build_name"`
	PackerBuilderType   *string           `mapstructure:"packer_builder_type" cty:"packer_builder_type" hcl:"packer_builder_type"`
	PackerCoreVersion   *string           `mapstructure:"packer_core_version" cty:"packer_core_version" hcl:"packer_core_version"`
	PackerDebug         *bool             `mapstructure:"packer_debug" cty:"packer_debug" hcl:"packer_debug"`
	PackerForce         *bool             `mapstructure:"packer_force" cty:"packer_force" hcl:"packer_force"`
	PackerOnError       *string           `mapstructure:"packer_on_error" cty:"packer_on_error" hcl:"packer_on_error"`
	PackerUserVars      map[string]string `mapstructure:"packer_user_variables" cty:"packer_user_variables" hcl:"packer_user_variables"`
	PackerSensitiveVars []string          `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables" hcl:"packer_sensitive_variables"`
	KsyunAccessKey      *string           `mapstructure:"access_key" required:"true" cty:"access_key" hcl:"access_key"`
	KsyunSecretKey      *string           `mapstructure:"secret_key" required:"true" cty:"secret_key" hcl:"secret_key"`
	KsyunRegion         *string           `mapstructure:"region" required:"true" cty:"region" hcl:"region"`
	MostRecent          *bool             `mapstructure:"most_recent" cty:"most_recent" hcl:"most_recent"`
	ImageSource         *string           `mapstructure:"image_source" cty:"image_source" hcl:"image_source"`
	NameRegex           *string           `mapstructure:"name_regex" cty:"name_regex" hcl:"name_regex"`
	Platform            *string           `mapstructure:"platform" cty:"platform" hcl:"platform"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"packer_build_name":          &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":        &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_core_version":        &hcldec.AttrSpec{Name: "packer_core_version", Type: cty.String, Required: false},
		"packer_debug":               &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":               &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":            &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":      &hcldec.AttrSpec{Name: "packer_user_variables", Type: cty.Map(cty.String), Required: false},
		"packer_sensitive_variables": &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
		"access_key":                 &hcldec.AttrSpec{Name: "access_key", Type: cty.String, Required: false},
		"secret_key":                 &hcldec.AttrSpec{Name: "secret_key", Type: cty.String, Required: false},
		"region":                     &hcldec.AttrSpec{Name: "region", Type: cty.String, Required: false},
		"most_recent":                &hcldec.AttrSpec{Name: "most_recent", Type: cty.Bool, Required: false},
		"image_source":               &hcldec.AttrSpec{Name: "image_source", Type: cty.String, Required: false},
		"name_regex":                 &hcldec.AttrSpec{Name: "name_regex", Type: cty.String, Required: false},
		"platform":                   &hcldec.AttrSpec{Name: "platform", Type: cty.String, Required: false},
	}
	return s
}

// FlatDatasourceOutput is an auto-generated flat version of DatasourceOutput.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatDatasourceOutput struct {
	ID           *string `mapstructure:"id" cty:"id" hcl:"id"`
	Name         *string `mapstructure:"name" cty:"name" hcl:"name"`
	CreationDate *string `mapstructure:"creation_date" cty:"creation_date" hcl:"creation_date"`
	Platform     *string `mapstructure:"platform" cty:"platform" hcl:"platform"`
	ImageSource  *string `mapstructure:"image_source" cty:"image_source" hcl:"image_source"`
}

// FlatMapstructure returns a new FlatDatasourceOutput.
// FlatDatasourceOutput is an auto-generated flat version of DatasourceOutput.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*DatasourceOutput) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatDatasourceOutput)
}

// HCL2Spec returns the hcl spec of a DatasourceOutput.
// This spec is used by HCL to read the fields of DatasourceOutput.
// The decoded values from this spec will then be applied to a FlatDatasourceOutput.
func (*FlatDatasourceOutput) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"id":            &hcldec.AttrSpec{Name: "id", Type: cty.String, Required: false},
		"name":          &hcldec.AttrSpec{Name: "name", Type: cty.String, Required: false},
		"creation_date": &hcldec.AttrSpec{Name: "creation_date", Type: cty.String, Required: false},
		"platform":      &hcldec.AttrSpec{Name: "platform", Type: cty.String, Required: false},
		"image_source":  &hcldec.AttrSpec{Name: "image_source", Type: cty.String, Required: false},
	}
	return s
}
