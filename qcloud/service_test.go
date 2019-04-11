package qcloud

import (
	"testing"
	"fmt"
)

func TestServiceBmlb(t *testing.T)  {
	service := NewService(Bmlb, BJ)
	params := make(map[string]interface{})
	params["loadBalancerType"] = "internal"
	params["offset"] = 0
	params["limit"] = 10

	resp := service.call(DescribeBmLoadBalancers, params)
	fmt.Println(resp)
}

func TestServiceBmvpc(t *testing.T)  {
	service := NewService(Bmvpc, BJ)
	params := make(map[string]interface{})
	resp := service.call(DescribeBmVpcEx, params)
	fmt.Println(resp)
}
