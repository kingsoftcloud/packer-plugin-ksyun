package kec

import (
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/kingsoftcloud/packer-plugin-ksyun/builder"
	"time"
)

type ClientKecWrapper struct {
	*ksyun.ClientWrapper
	KecClient *kec.Kec
}

const (
	defaultKecInstanceName = "ksyun_packer_vm"
	defaultKecInstanceType = "I1.1A"
	defaultKecSshUserName  = "root"
	defaultKecChargeType   = "Daily"
)

func (c *ClientKecWrapper) WaitKecInstanceStatus(stateBag multistep.StateBag, instanceId string, projectId string, status string) (*map[string]interface{}, error) {
	return c.WaitResource(&ksyun.WaitResourceParam{
		RequestResource: func() (*map[string]interface{}, error) {
			queryKec := make(map[string]interface{})
			queryKec["InstanceId.1"] = instanceId
			queryKec["ProjectId.1"] = projectId
			return c.KecClient.DescribeInstances(&queryKec)
		},
		ProcessRequest: func(resp *map[string]interface{}, err error) ksyun.ProcessRequestResult {
			if err != nil {
				return ksyun.RequestResourceRetry
			}
			kecId := ksyun.GetSdkValue(stateBag, "InstancesSet.0.InstanceId", *resp)
			if kecId == nil {
				return ksyun.RequestResourceRetry
			}
			kecState := ksyun.GetSdkValue(stateBag, "InstancesSet.0.InstanceState.Name", *resp).(string)
			if kecState == status {
				if stateBag.Get("PrivateIp") == nil {
					privateIp := ksyun.GetSdkValue(stateBag, "InstancesSet.0.NetworkInterfaceSet.0.PrivateIpAddress", *resp).(string)
					stateBag.Put("PrivateIp", privateIp)
				}

				return ksyun.RequestResourceSuccess
			}
			return ksyun.RequestResourceRetry
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    360,
	})
}

func (c *ClientKecWrapper) WaitKecImageStatus(stateBag multistep.StateBag, imageId string, status string) (*map[string]interface{}, error) {
	return c.WaitResource(&ksyun.WaitResourceParam{
		RequestResource: func() (*map[string]interface{}, error) {
			queryImage := make(map[string]interface{})
			queryImage["ImageId"] = imageId
			return c.KecClient.DescribeImages(&queryImage)
		},
		ProcessRequest: func(resp *map[string]interface{}, err error) ksyun.ProcessRequestResult {
			if err != nil {
				return ksyun.RequestResourceRetry
			}
			id := ksyun.GetSdkValue(stateBag, "ImagesSet.0.ImageId", *resp)
			if id == nil {
				return ksyun.RequestResourceRetry
			}
			state := ksyun.GetSdkValue(stateBag, "ImagesSet.0.ImageState", *resp).(string)
			if state == status {
				return ksyun.RequestResourceSuccess
			}
			return ksyun.RequestResourceRetry
		},
		RetryInterval: 30 * time.Second,
		RetryTimes:    360,
	})
}
