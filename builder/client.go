package ksyun

import (
	"fmt"
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"time"
)

const (
	defaultRetryInterval = 5 * time.Second
	defaultRetryTimes    = 12
)
const (
	defaultProjectId = "0"
)

const (
	defaultVpcName           = "ksyun_packer_vpc"
	defaultVpcCidr           = "172.20.0.0/16"
	defaultSubnetName        = "ksyun_packer_subnet"
	defaultSubnetCidr        = "172.20.1.0/24"
	defaultSubnetType        = "Normal"
	defaultSecurityGroupName = "ksyun_packer_security_group"
)

type ClientWrapper struct {
	SksClient *sks.Sks
	EipClient *eip.Eip
	VpcClient *vpc.Vpc
	//KecClient *kec.Kec
	//EpcClient *epc.Epc
}

type ProcessRequestResult struct {
	Complete  bool
	StopRetry bool
}

var (
	RequestResourceSuccess = ProcessRequestResult{
		Complete:  true,
		StopRetry: true,
	}

	RequestResourceRetry = ProcessRequestResult{
		Complete:  false,
		StopRetry: false,
	}

	//RequestResourceStop = ProcessRequestResult{
	//	Complete:  false,
	//	StopRetry: true,
	//}
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
		if processRequestResult.Complete {
			return response, nil
		}
		if processRequestResult.StopRetry {
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
			networkInterfaces := GetSdkValue(stateBag, "NetworkInterfaceSet", *resp).([]interface{})
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
			networkInterfaces := GetSdkValue(stateBag, "NetworkInterfaceSet", *resp).([]interface{})
			if len(networkInterfaces) == 0 {
				return RequestResourceSuccess
			}
			return RequestResourceRetry
		},
		RetryInterval: 10 * time.Second,
		RetryTimes:    360,
	})
}
