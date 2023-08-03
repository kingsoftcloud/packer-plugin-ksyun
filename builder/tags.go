package ksyun

import (
	"fmt"
	"strconv"

	"github.com/KscSDK/ksc-sdk-go/service/tagv2"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

const (
	ResourceTypeKec       = "instance"
	ResourceTypeEip       = "eip"
	ResourceTypeImage     = "image"
	ResourceTypeBareHost  = "epc-host"
	ResourceTypeBareImage = "epc-image"
)

var ResourceTypeList = []string{ResourceTypeImage, ResourceTypeKec, ResourceTypeEip}

type TagMap map[string]string

type Tag struct {
	Id         int    `mapstructure:"Id"`
	Key        string `mapstructure:"Key"`
	Value      string `mapstructure:"Value"`
	CreateTime string `mapstructure:"CreateTime"`
	CanDelete  int    `mapstructure:"CanDelete"`
	IsBillTag  int    `mapstructure:"IsBillTag"`
}

type Tags []*Tag

func (t Tags) Report(ui packersdk.Ui) {
	for _, tag := range t {
		ui.Message(fmt.Sprintf("Adding tag: \"%s\": \"%s\"",
			tag.Key, tag.Value))
	}
}

func (t TagMap) KsyunTags() Tags {
	var tags = Tags{}
	for k, v := range t {
		tags = append(tags, &Tag{Key: k, Value: v})
	}
	return tags
}

func (t Tags) GetTagsParams(rsType, rsUuid string) (map[string]interface{}, error) {
	if !StringInSlice(rsType, ResourceTypeList, false) {
		return nil, fmt.Errorf("resource type is not empty, specify it")
	}
	if !IsUuid(rsUuid) {
		return nil, fmt.Errorf("resouce uuid is invalid")
	}
	m := make(map[string]interface{})
	for i, tag := range t {
		m["Tag_"+strconv.Itoa(i+1)+"_Key"] = tag.Key
		m["Tag_"+strconv.Itoa(i+1)+"_Value"] = tag.Value
	}
	m["ResourceType"] = rsType
	rpTagMap := map[string]interface{}{
		"ResourceUuids": rsUuid,
	}
	m["ReplaceTags"] = []interface{}{rpTagMap}
	return m, nil
}

func (t Tags) Pinning(rsType, rsUuid string, client *tagv2.Tagv2) error {
	tagsParams, err := t.GetTagsParams(rsType, rsUuid)
	if err != nil {
		return err
	}

	_, err = client.ReplaceResourcesTags(&tagsParams)
	if err != nil {
		return err
	}
	return nil
}
