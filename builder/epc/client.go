package epc

import (
	"github.com/KscSDK/ksc-sdk-go/service/eip"
	"github.com/KscSDK/ksc-sdk-go/service/epc"
	"github.com/KscSDK/ksc-sdk-go/service/sks"
	"github.com/KscSDK/ksc-sdk-go/service/vpc"
)

type ClientWrapper struct {
	EpcClient *epc.Epc
	SksClient *sks.Sks
	EipClient *eip.Eip
	VpcClient *vpc.Vpc
}
