package ksyun

import (
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"reflect"
	"strconv"
	"strings"
)

func SSHHost(c communicator.Config) func(multistep.StateBag) (string, error) {
	return func(stateBag multistep.StateBag) (string, error) {
		publicIp := stateBag.Get("publicIp").(string)
		return publicIp, nil
	}
}

func Halt(stateBag multistep.StateBag, err error, prefix string) multistep.StepAction {
	ui := stateBag.Get("ui").(packersdk.Ui)

	if prefix != "" {
		err = fmt.Errorf("%s: %s", prefix, err)
	}

	stateBag.Put("error", err)
	ui.Error(err.Error())
	return multistep.ActionHalt
}

func GetSdkValue(stateBag multistep.StateBag, keyPattern string, obj interface{}) interface{} {
	keys := strings.Split(keyPattern, ".")
	root := obj
	for index, k := range keys {
		if reflect.ValueOf(root).Kind() == reflect.Map {
			root = root.(map[string]interface{})[k]
			if root == nil {
				return root
			}

		} else if reflect.ValueOf(root).Kind() == reflect.Slice {
			i, err := strconv.Atoi(k)
			if err != nil {
				Halt(stateBag, err, fmt.Sprintf("keyPattern %s index %d must number", keyPattern, index))
			}
			if len(root.([]interface{})) <= i {
				return nil
			}
			root = root.([]interface{})[i]
		}
	}
	return root
}

func GetCidrIpRange(cidr string) (string, string, string) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg3MinIp, seg3MaxIp := GetIpSeg3Range(ipSegs, maskLen)
	seg4MinIp, seg4MaxIp := GetIpSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."

	return ipPrefix + strconv.Itoa(seg3MinIp) + "." + strconv.Itoa(seg4MinIp),
		ipPrefix + strconv.Itoa(seg3MinIp) + "." + strconv.Itoa(seg4MinIp+1),
		ipPrefix + strconv.Itoa(seg3MaxIp) + "." + strconv.Itoa(seg4MaxIp-2)
}

func GetCidrHostNum(maskLen int) uint {
	cidrIpNum := uint(0)
	var i uint = uint(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIpNum += 1 << i
	}
	return cidrIpNum
}

func GetCidrIpMask(maskLen int) string {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-maskLen)
	fmt.Println(fmt.Sprintf("%b \n", cidrMask))
	//计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))

	return fmt.Sprint(cidrMaskSeg1) + "." + fmt.Sprint(cidrMaskSeg2) + "." + fmt.Sprint(cidrMaskSeg3) + "." + fmt.Sprint(cidrMaskSeg4)
}

func GetIpSeg3Range(ipSegs []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIp, _ := strconv.Atoi(ipSegs[2])
		return segIp, segIp
	}
	ipSeg, _ := strconv.Atoi(ipSegs[2])
	return GetIpSegRange(uint8(ipSeg), uint8(24-maskLen))
}

func GetIpSeg4Range(ipSegs []string, maskLen int) (int, int) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	segMinIp, segMaxIp := GetIpSegRange(uint8(ipSeg), uint8(32-maskLen))
	return segMinIp + 1, segMaxIp - 2
}

func GetIpSegRange(userSegIp, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << offset
	segMinIp := netSegIp & userSegIp
	segMaxIp := userSegIp&(255<<offset) | ^(255 << offset)
	return int(segMinIp), int(segMaxIp)
}
