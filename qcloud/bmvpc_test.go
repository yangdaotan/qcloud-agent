package qcloud

import (
	"testing"
	"fmt"
)

func TestBmvpcServiceCall(t *testing.T)  {
	bmvpcService := &BmvpcService{
		Service: NewService(Bmvpc, BJ),
	}

	params := make(map[string]interface{})

	fmt.Println(bmvpcService.call(DescribeBmVpcEx, params))
}

func TestGetMatchVpc(t *testing.T)  {
	bmvpcService := &BmvpcService{
		Service: NewService(Bmvpc, BJ),
	}

	fmt.Println(bmvpcService.getMatchVpc("10.201.95.2"))
}

func TestDescribeBmSubnetEx(t *testing.T)  {
	bmvpcService := &BmvpcService{
		Service: NewService(Bmvpc, BJ),
	}

	fmt.Println(bmvpcService.getAvailableSubnet("vpc-08k4chpj"))
}