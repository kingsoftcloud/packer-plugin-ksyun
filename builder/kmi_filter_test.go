package ksyun

import (
	"os"
	"testing"

	"github.com/KscSDK/ksc-sdk-go/ksc"
	"github.com/KscSDK/ksc-sdk-go/ksc/utils"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
)

func TestKmiFilterOptions_GetFilteredImage(t *testing.T) {

	kecConn := GetKecClient("")
	kfo := KmiFilterOptions{
		ImageSource: "system",
		Platform:    "centos-7.5",
		MostRecent:  true,
	}
	params := map[string]interface{}{
		"ImageId": "IMG-12112384-c3d3-4d42-8882-58234825ba1c",
	}
	data, err := kfo.GetFilteredImage(&params, kecConn)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func GetKecClient(region string) *kec.Kec {
	ak := os.Getenv("KSYUN_ACCESS_KEY")
	sk := os.Getenv("KSYUN_SECRET_KEY")
	if region == "" {
		region = "cn-beijing-6"
	}
	cli := ksc.NewClient(ak, sk)
	cfg := &ksc.Config{
		Region: &region,
	}
	url := &utils.UrlInfo{
		UseSSL: false,
		Locate: false,
	}

	return kec.SdkNew(cli, cfg, url)
}
