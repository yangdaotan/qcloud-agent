package qcloud

import (
	"fmt"
	"errors"
)

type CnsService struct {
	Service *Service
}

func (cnsService *CnsService) Call(action Action, params map[string]interface{}) string {
	return cnsService.Service.call(action, params)
}

/**
	rockemtq域名：mq.$env.bike.io
	kafka域名：kafka.$env.bike.io
	zk域名：zk.$env.bike.io
 */

func (cnsService *CnsService) GetRecord(domain, subDomain string) (int64, error) {
	params := make(map[string]interface{})
	params["domain"] = "bike.io"
	params["subDomain"] = subDomain

	resp := cnsService.Call(RecordList, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return -1, errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	records := dat["data"].(map[string]interface{})["records"].([]interface{})


	if records == nil {
		return -1, nil
	}

	fmt.Printf("%v\n", records)

	return records[0].(map[string]interface{})["id"].(int64), nil
}