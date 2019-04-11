package qcloud

import (
	"testing"
	"fmt"
)

func TestGetDevice(t *testing.T )  {
	bmService := &BmService{
		Service:NewService(Bm, BJ),
	}

	ips := make(map[string]interface{})
	ips["lanIps.0"] = "10.50.0.55"
	ips["lanIps.1"] = "10.50.0.56"

	fmt.Println(bmService.GetDevices(ips))
}
