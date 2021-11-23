package epc

import (
	"github.com/KscSDK/ksc-sdk-go/service/epc"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	ksyun "github.com/kingsoftcloud/packer-plugin-ksyun/builder"
	"time"
)

type ClientEpcWrapper struct {
	*ksyun.ClientWrapper
	EpcClient *epc.Epc
}

const (
	defaultEpcInstanceName = "ksyun_packer_bare_metal"
	defaultEpcChargeType   = "Daily"
	defaultEpcImageName    = "packer_bare_metal_image"
)

func (c *ClientEpcWrapper) WaitEpcInstanceStatus(stateBag multistep.StateBag, instanceId string, projectId string, status string) (*map[string]interface{}, error) {
	return c.WaitResource(&ksyun.WaitResourceParam{
		RequestResource: func() (*map[string]interface{}, error) {
			queryKec := make(map[string]interface{})
			queryKec["HostId.1"] = instanceId
			queryKec["ProjectId.1"] = projectId
			return c.EpcClient.DescribeEpcs(&queryKec)
		},
		ProcessRequest: func(resp *map[string]interface{}, err error) ksyun.ProcessRequestResult {
			if err != nil {
				return ksyun.RequestResourceRetry
			}
			epcId := ksyun.GetSdkValue(stateBag, "HostSet.0.HostId", *resp)
			if epcId == nil {
				return ksyun.RequestResourceRetry
			}
			epcState := ksyun.GetSdkValue(stateBag, "HostSet.0.HostStatus", *resp).(string)
			if epcState == status {
				if stateBag.Get("PrivateIp") == nil {
					privateIps := ksyun.GetSdkValue(stateBag, "HostSet.0.NetworkInterfaceAttributeSet", *resp).([]interface{})
					for _, p := range privateIps {
						networkInterfaceType := ksyun.GetSdkValue(stateBag, "NetworkInterfaceType", p).(string)
						if networkInterfaceType == "primary" {
							privateIp := ksyun.GetSdkValue(stateBag, "PrivateIpAddress", p).(string)
							stateBag.Put("PrivateIp", privateIp)
							break
						}
					}

				}

				return ksyun.RequestResourceSuccess
			}
			return ksyun.RequestResourceRetry
		},
		RetryInterval: 60 * time.Second,
		RetryTimes:    360,
	})
}
