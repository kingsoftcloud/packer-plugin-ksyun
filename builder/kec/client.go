package kec

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"time"
)

type ClientWrapper struct {
	KecClient *kec.Kec
	SksClient *sks.Sks
	EipClient *eip.Eip
	VpcClient *vpc.Vpc
}

const (
	defaultRetryInterval = 5 * time.Second
	defaultRetryTimes    = 12
)

const (
	defaultKecInstanceName = "kingcloud_packer_vm"
	defaultKecInstanceType = "I1.1A"
	defaultKecSshUserName  = "root"
	defaultKecChargeType   = "Daily"
)

const (
	defaultProjectId = "0"
)

const (
	defaultVpcName           = "kingcloud_packer_vpc"
	defaultVpcCidr           = "172.20.0.0/16"
	defaultSubnetName        = "kingcloud_packer_subnet"
	defaultSubnetCidr        = "172.20.1.0/24"
	EnableSubnetType         = "Normal"
	defaultSecurityGroupName = "kingcloud_packer_security_group"
)

const (
	defaultSSHKeyName = "kingcloud_packer_ssh_key"
)

type ProcessRequestResult struct {
	complete  bool
	stopRetry bool
}

var (
	RequestResourceSuccess = ProcessRequestResult{
		complete:  true,
		stopRetry: true,
	}

	RequestResourceRetry = ProcessRequestResult{
		complete:  false,
		stopRetry: false,
	}

	RequestResourceStop = ProcessRequestResult{
		complete:  false,
		stopRetry: true,
	}
)

type WaitResourceParam struct {
	RequestResource func() (*map[string]interface{}, error)
	ProcessRequest  func(*map[string]interface{}, error) ProcessRequestResult
	RetryInterval   time.Duration
	RetryTimes      int
}

func (c *ClientWrapper) WaitResource(param *WaitResourceParam) (*map[string]interface{}, error) {
	if param.RetryTimes <= 0 {
		param.RetryTimes = defaultRetryTimes
	}

	if param.RetryInterval <= 0 {
		param.RetryInterval = defaultRetryInterval
	}

	var lastResponse *map[string]interface{}
	var lastError error

	for i := 0; ; i++ {

		if i >= param.RetryTimes {
			break
		}

		response, err := param.RequestResource()
		lastResponse = response
		lastError = err

		processRequestResult := param.ProcessRequest(lastResponse, lastError)
		if processRequestResult.complete {
			return response, nil
		}
		if processRequestResult.stopRetry {
			return response, err
		}

		time.Sleep(param.RetryInterval)
	}

	if lastError == nil {
		lastError = fmt.Errorf("<no error>")
	}

	return lastResponse, fmt.Errorf("wait failed after %d times retry with %d seconds retry interval: %s",
		param.RetryTimes, int(param.RetryInterval.Seconds()), lastError)
}

func (c *ClientWrapper) WaitKecInstanceStatus(stateBag multistep.StateBag, instanceId string, projectId string, status string) (*map[string]interface{}, error) {
	return c.WaitResource(&WaitResourceParam{
		RequestResource: func() (*map[string]interface{}, error) {
			queryKec := make(map[string]interface{})
			queryKec["InstanceId.1"] = instanceId
			queryKec["ProjectId.1"] = projectId
			return c.KecClient.DescribeInstances(&queryKec)
		},
		ProcessRequest: func(resp *map[string]interface{}, err error) ProcessRequestResult {
			if err != nil {
				return RequestResourceRetry
			}
			kecId := getSdkValue(stateBag, "InstancesSet.0.InstanceId", *resp)
			if kecId == nil {
				return RequestResourceRetry
			}
			kecState := getSdkValue(stateBag, "InstancesSet.0.InstanceState.Name", *resp).(string)
			if kecState == status {
				return RequestResourceSuccess
			}
			return RequestResourceRetry
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    360,
	})
}

func (c *ClientWrapper) WaitKecImageStatus(stateBag multistep.StateBag, imageId string, status string) (*map[string]interface{}, error) {
	return c.WaitResource(&WaitResourceParam{
		RequestResource: func() (*map[string]interface{}, error) {
			queryImage := make(map[string]interface{})
			queryImage["ImageId"] = imageId
			return c.KecClient.DescribeImages(&queryImage)
		},
		ProcessRequest: func(resp *map[string]interface{}, err error) ProcessRequestResult {
			if err != nil {
				return RequestResourceRetry
			}
			id := getSdkValue(stateBag, "ImagesSet.0.ImageId", *resp)
			if id == nil {
				return RequestResourceRetry
			}
			state := getSdkValue(stateBag, "ImagesSet.0.ImageState", *resp).(string)
			if state == status {
				return RequestResourceSuccess
			}
			return RequestResourceRetry
		},
		RetryInterval: 30 * time.Second,
		RetryTimes:    360,
	})
}

func (c *ClientWrapper) WaitSecurityGroupClean(stateBag multistep.StateBag, securityGroupId string) (*map[string]interface{}, error) {
	return c.WaitResource(&WaitResourceParam{
		RequestResource: func() (*map[string]interface{}, error) {
			queryNetworkInterfaces := make(map[string]interface{})
			queryNetworkInterfaces["Filter.1.Name"] = "securitygroup-id"
			queryNetworkInterfaces["Filter.1.Value.1"] = securityGroupId
			return c.VpcClient.DescribeNetworkInterfaces(&queryNetworkInterfaces)
		},
		ProcessRequest: func(resp *map[string]interface{}, err error) ProcessRequestResult {
			if err != nil {
				return RequestResourceRetry
			}
			networkInterfaces := getSdkValue(stateBag, "NetworkInterfaceSet", *resp).([]interface{})
			if len(networkInterfaces) == 0 {
				return RequestResourceSuccess
			}
			return RequestResourceRetry
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    360,
	})
}

func (c *ClientWrapper) WaitSubnetClean(stateBag multistep.StateBag, subnetId string) (*map[string]interface{}, error) {
	return c.WaitResource(&WaitResourceParam{
		RequestResource: func() (*map[string]interface{}, error) {
			queryNetworkInterfaces := make(map[string]interface{})
			queryNetworkInterfaces["Filter.1.Name"] = "subnet-id"
			queryNetworkInterfaces["Filter.1.Value.1"] = subnetId
			return c.VpcClient.DescribeNetworkInterfaces(&queryNetworkInterfaces)
		},
		ProcessRequest: func(resp *map[string]interface{}, err error) ProcessRequestResult {
			if err != nil {
				return RequestResourceRetry
			}
			networkInterfaces := getSdkValue(stateBag, "NetworkInterfaceSet", *resp).([]interface{})
			if len(networkInterfaces) == 0 {
				return RequestResourceSuccess
			}
			return RequestResourceRetry
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    360,
	})
}
