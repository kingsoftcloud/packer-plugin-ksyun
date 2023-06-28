package kec

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
)

const (
	ImageActionShare = "share"
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

		// share image
		err = s.ImageShare(imageId, stateBag)
		if err != nil {
			return ksyun.Halt(stateBag, err, "Error sharing kec image")
		}
		if s.KsyunImageConfig.KsyunImageWarmUp {
			if err := s.ImageWarmup(imageId, stateBag); err != nil {
				return ksyun.Halt(stateBag, err, "Error warming up image")
			}

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
			// return ksyun.Halt(stateBag, err, "error copying images")
		}
		ui.Message(fmt.Sprintf("copy image to %s", region))

	}
	return nil

}

// ImageShare to deal with an image share to other account
func (s *stepCreateKsyunImage) ImageShare(imageId string, stateBag multistep.StateBag) error {
	shareAccounts := s.KsyunImageConfig.KsyunImageShareAccounts

	if len(shareAccounts) == 0 {
		return nil
	}

	ui := stateBag.Get("ui").(packersdk.Ui)

	client := stateBag.Get("client").(*ClientKecWrapper)
	for _, shareAccount := range shareAccounts {
		params := map[string]interface{}{
			"ImageId":     imageId,
			"AccountId.1": shareAccount,
		}

		params["Permission"] = ImageActionShare
		_, err := client.KecClient.ModifyImageSharePermission(&params)
		if err != nil {
			return err
			// return ksyun.Halt(stateBag, err, "error copying images")
		}
		ui.Message(fmt.Sprintf("Image share to %s", shareAccount))

	}
	return nil
}

// ImageWarmup set this image start to warmup-start
func (s *stepCreateKsyunImage) ImageWarmup(imageId string, stateBag multistep.StateBag) error {
	ui := stateBag.Get("ui").(packersdk.Ui)
	client := stateBag.Get("client").(*ClientKecWrapper)

	if _, err := client.WaitKecImageStatus(stateBag, imageId, "active"); err != nil {
		return err
	}

	params := map[string]interface{}{
		"ImageId.1": imageId,
	}

	// EnableImageCaching is an async call, so we have to wait its state changed
	if _, err := client.KecClient.EnableImageCaching(&params); err != nil {
		return err
	}
	// ui.Say("Waiting image warming up")
	// waiting image warm-up state syncing
	time.Sleep(10 * time.Second)
	// TODO: 在action=describeImages中若不添加header: X-KSC-SOURCE=kec 则无法获取到images的warm-up状态
	if _, err := client.WaitKecImageStatus(stateBag, imageId, "active"); err != nil {
		return err
	}
	ui.Message("Image set warm-up successfully")

	return nil
}

func (s *stepCreateKsyunImage) Cleanup(stateBag multistep.StateBag) {

}
