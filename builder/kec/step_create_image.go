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

	if s.KsyunImageConfig.KsyunImageType != "" {
		createImage["Type"] = s.KsyunImageConfig.KsyunImageType
	}

	// 判断是否创建整机镜像
	if !s.KsyunImageConfig.KsyunImageIgnoreDataDisks {
		dataDisksSrc := reflect.ValueOf(stateBag.Get("DataDisks"))
		if dataDisksSrc.Kind() == reflect.Slice && dataDisksSrc.Len() > 0 {
			for i := 0; i < dataDisksSrc.Len(); i++ {
				log.Println("dataDisksSrc:", i)
				ele := dataDisksSrc.Index(i).Elem()
				createImage[fmt.Sprintf("DataDiskIds.%d", i+1)] = ele.MapIndex(reflect.ValueOf("DiskId")).Elem().String()
			}
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

		// copy image
		err = s.ImageCopy(imageId, stateBag)
		if err != nil {
			return ksyun.Halt(stateBag, err, "Error copying kec image")
		}
	}

	return multistep.ActionContinue
}

func (s *stepCreateKsyunImage) ImageCopy(imageId string, stateBag multistep.StateBag) error {
	regions := s.KsyunImageConfig.KsyunImageCopyRegions

	if len(regions) == 0 {
		return nil
	}

	names := s.KsyunImageConfig.KsyunImageCopyNames
	ui := stateBag.Get("ui").(packersdk.Ui)

	client := stateBag.Get("client").(*ClientKecWrapper)
	for idx, region := range regions {
		name := ""
		if idx < len(names) {
			name = names[idx]
		}
		params := map[string]interface{}{
			"ImageId.1":           imageId,
			"DestinationRegion.1": region,
		}
		if name != "" {
			params["DestinationImageName"] = name
		}
		// 目前返回值没有新镜像的id，不能查进度
		_, err := client.KecClient.CopyImage(&params)
		if err != nil {
			return err
			//return ksyun.Halt(stateBag, err, "error copying images")
		}
		ui.Message(fmt.Sprintf("copy image to %s", region))

	}
	return nil

}

func (s *stepCreateKsyunImage) Cleanup(stateBag multistep.StateBag) {

}
