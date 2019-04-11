package qcloud

import (
	"fmt"
	"net"
)

type BmvpcService struct {
	Service *Service
}

type Vpc struct {
	vpcId int64
	unVpcId string
	vpcName string
	cidrBlock string
	subnetNum int64
}


func (bmvpcService *BmvpcService) getMatchVpc(nodeIP string) *Vpc {
	vpcs := bmvpcService.call(DescribeBmVpcEx, nil)

	error, data := CheckSuccess(vpcs)
	if error != nil {
		fmt.Printf("curl qcould error, %v", error)
	}

	for _, vpc := range data.([]interface{}) {
		vpcMap := vpc.(map[string]interface{})
		cidrBlock := vpcMap["cidrBlock"]
		_, ipNet, _ := net.ParseCIDR(cidrBlock.(string))
		if ipNet.Contains(net.ParseIP(nodeIP)) {
			return &Vpc{
				vpcId:vpcMap["vpcId"].(int64),
				unVpcId:vpcMap["unVpcId"].(string),
				vpcName:vpcMap["vpcName"].(string),
				cidrBlock:vpcMap["cidrBlock"].(string),
				subnetNum: vpcMap["subnetNum"].(int64),
			}
		}
	}

	return nil
}

//See :https://cloud.tencent.com/document/product/386/6648
func (bmvpcService *BmvpcService) getAvailableSubnet(unVpcId string) string {
	params := make(map[string]interface{})
	params["unVpcId"] = unVpcId

	subnets := bmvpcService.call(DescribeBmSubnetEx, params)
	error, data := CheckSuccess(subnets)
	if error != nil {
		fmt.Printf("curl qcould error, %v", error)
	}

	var maxAvailableIpNum int64 = 0
	var unSubnetId string
	for _, subnet := range data.([]interface{}) {
		subnetMap := subnet.(map[string]interface{})
		availableIpNum := subnetMap["availableIpNum"].(int64)
		if maxAvailableIpNum < availableIpNum {
			maxAvailableIpNum = availableIpNum
			unSubnetId = subnetMap["unSubnetId"].(string)
		}
	}

	if maxAvailableIpNum > 0 {
		return unSubnetId
	}

	return ""
}

func (bmvpcService *BmvpcService) call(action Action, params map[string]interface{}) string {
	return bmvpcService.Service.call(action, params)
}

