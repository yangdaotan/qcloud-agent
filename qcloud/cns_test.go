package qcloud

import (
	"testing"
	"fmt"
)

func TestDomainList(t *testing.T)  {
	cnsService := &CnsService{
		Service: NewService(Cns, BJ),
	}

	fmt.Println(cnsService.Call(DomainList, nil))
}

/**
       		{
                "id": 393541268,
                "ttl": 600, 
                "value": "10.201.194.3",
                "enabled": 1,
                "status": "enabled",
                "updated_on": "2018-11-27 16:37:53",
                "q_project_id": 0,
                "name": "kafka.longines",
                "line": "默认",
                "line_id": "0",
                "type": "A",
                "remark": "",
                "mx": 0
            },
 */
func TestRecordList(t *testing.T)  {
	cnsService := &CnsService{
		Service: NewService(Cns, BJ),
	}

	params := make(map[string]interface{})
	params["domain"] = "bike.io"
	params["subDomain"] = "rocketmq.beta"

	fmt.Println(cnsService.Call(RecordList, params))
}


/*
	{
		"code":0,
		"message":"",
		"codeDesc":"Success",
		"data":{
			"record":{
				"id":"395983060",
				"name":"kafka.mudun",
				"status":"enabled",
				"weight":null
			}
		}
	}
*/

//See: https://cloud.tencent.com/document/product/302/8516
func TestRecordCreate(t *testing.T)  {
	cnsService := &CnsService{
		Service: NewService(Cns, BJ),
	}

	params := make(map[string]interface{})
	params["domain"] = "bike.io"
	params["subDomain"] = "omega.mq"
	params["recordType"] = "A"
	params["recordLine"] = "默认"
	// 10.201.95.2的lb地址为10.201.116.2
	// 10.201.116.2:9092 -> 10.201.95.2:30003
	params["value"] = "10.190.3.42"
	params["ttl"] = 10

	fmt.Println(cnsService.Call(RecordCreate, params))
}

func TestRecordDelete(t *testing.T)  {
	cnsService := &CnsService{
		Service: NewService(Cns, BJ),
	}

	params := make(map[string]interface{})
	params["domain"] = "bike.io"
	params["recordId"] = 409702464

	fmt.Println(cnsService.Call(RecordDelete, params))
}

func TestRecordGet(t *testing.T)  {
	cnsService := &CnsService{
		Service: NewService(Cns, BJ),
	}

	fmt.Println(cnsService.GetRecord("bike.io", "mq.ydt"))
}
