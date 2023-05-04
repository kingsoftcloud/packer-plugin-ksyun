package kec

import (
	"context"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
	"log"
	"reflect"
)

type stepCreateKsyunImage struct {
	KsyunRunConfig   *KsyunKecRunConfig
	KsyunImageConfig *KsyunImageConfig
}

func (s *stepCreateKsyunImage) Run(ctx context.Context, stateBag multistep.StateBag) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientKecWrapper)
	instanceId := stateBag.Get("InstanceId").(string)
	ui.Say("Creating Ksyun Kec Image ")
	createImage := make(map[string]interface{})
	createImage["InstanceId"] = instanceId
	createImage["Name"] = s.KsyunImageConfig.KsyunImageName
	//log.Println("11111", s.KsyunImageConfig, s.KsyunRunConfig)
	if s.KsyunImageConfig.KsyunImageType != "" {
		createImage["Type"] = s.KsyunImageConfig.KsyunImageType
	}

	dataDisksSrc := reflect.ValueOf(stateBag.Get("DataDisks"))
	if dataDisksSrc.Kind() == reflect.Slice && dataDisksSrc.Len() > 0 {
		for i := 0; i < dataDisksSrc.Len(); i++ {
			log.Println("dataDisksSrc:", i)
			ele := dataDisksSrc.Index(i).Elem()
			createImage[fmt.Sprintf("DataDiskIds.%d", i+1)] = ele.MapIndex(reflect.ValueOf("DiskId")).Elem().String()
		}
	}
	log.Println("createImage:", createImage)

	resp, errorCreate := client.KecClient.CreateImage(&createImage)
	if errorCreate != nil {
		return ksyun.Halt(stateBag, errorCreate, "Error creating  kec image")
	}
	if resp != nil {
		ui.Say("Waiting Ksyun Kec Image active")
		imageId := ksyun.GetSdkValue(stateBag, "ImageId", *resp).(string)
		_, err := client.WaitKecImageStatus(stateBag, imageId, "active")
		if err != nil {
			return ksyun.Halt(stateBag, err, "Error waiting  kec image active")
		}
		stateBag.Put("TargetImageId", imageId)
	}
	return multistep.ActionContinue
}

func (s *stepCreateKsyunImage) Cleanup(stateBag multistep.StateBag) {

}
