package qcloud

import (
	"fmt"
)

type BmService struct {
	Service *Service
}

func (bmService *BmService) call(action Action, params map[string]interface{}) string {
	return bmService.Service.call(action, params)
}

// See: https://cloud.tencent.com/document/product/386/6728
func (bmService *BmService) GetDevices(lanIps map[string]interface{}) map[string]interface{} {
	devices := bmService.Service.call(DescribeDevice, lanIps)

	if devices == "" {
		return nil
	}

	error, data := CheckSuccess(devices)
	if error != nil || data == nil  {
		fmt.Printf("curl qcould error, %v", error)
		return nil
	}

	dataMap := data.(map[string]interface{})
	if dataMap == nil || dataMap["totalNum"].(int64) < 1 {
		return nil
	}

	instanceIds := make(map[string]interface{})
	for _, device := range dataMap["deviceList"].([]interface{}) {
		instanceIds[device.(map[string]interface{})["lanIp"].(string)] = device.(map[string]interface{})["instanceId"].(string)
	}

	return instanceIds
}
