package qcloud

import (
	"testing"
	"fmt"
	"time"
	"strings"
	"regexp"
)

// {"code":0,"message":"","codeDesc":"Success","loadBalancerIds":["lb-ny0xkppj"],"requestId":2915276}
func TestCreate(t *testing.T)  {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	params := make(map[string]interface{})
	params["nodeIp"] = "10.200.43.2"
	params["loadBalancerType"] = "internal"
	params["goodsNum"] = 1


	_, loadBalancerId := bmlbService.CreateLB(params)

	paramsT := make(map[string]interface{})
	paramsT["loadBalancerId"] = loadBalancerId
	paramsT["loadBalancerName"] = "k8s-mq"

	time.Sleep(time.Second * 5)

	fmt.Println(bmlbService.Call(ModifyBmLoadBalancerAttributes, paramsT))
}

func TestDescribeBmLoadBalancer(t *testing.T)  {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	params := make(map[string]interface{})
	params["loadBalancerType"] = "internal"
	params["limit"] = 50
	params["offset"] = 0


	lbs := make(map[string]string)
	for i := 0; i < 10; i = i+1 {
		params["loadBalancerType"] = "internal"
		params["limit"] = 50
		params["offset"] = 50 * i
		resp := bmlbService.Call(DescribeBmLoadBalancers, params)

		var dat map[string]interface{}
		if err := Unmarshal([]byte(resp), &dat); err != nil {
			fmt.Println(err)
		}

		loadBalancerSet := dat["loadBalancerSet"].([]interface{})

		for _, lb := range loadBalancerSet {
			//fmt.Printf(lb.(map[string]interface{})["loadBalancerName"])
			lbs[lb.(map[string]interface{})["loadBalancerName"].(string)] = lb.(map[string]interface{})["loadBalancerId"].(string)
		}

		ClearMap(params)

	}
	//fmt.Printf("total count: %v\n", len(lbs))
	for k, _ := range lbs {
		//fmt.Printf("%s %s\n", k, v)
		if !strings.Contains(k,"baal") {
			delete(lbs, k)
		}

		if find, _ := regexp.Match("\\d$", []byte(k)); !find {
			delete(lbs, k)
		}

		if strings.Contains(k, "diablo") {
			delete(lbs, k)
		}

		if strings.Contains(k, "talent") {
			delete(lbs, k)
		}

		if strings.Contains(k, "yubai") {
			delete(lbs, k)
		}

		if strings.Contains(k, "st1") {
			delete(lbs, k)
		}
	}

	fmt.Printf("total baal count: %v\n", len(lbs))
	for k, v := range lbs {
		fmt.Printf("-------deleting lb: %s %s--------\n", k, v)
		ClearMap(params)
		params["loadBalancerId"] = v
		fmt.Println(bmlbService.Call(DeleteBmLoadBalancers, params))
		fmt.Printf("--------------------------------\n\n")
		time.Sleep(10*time.Second)
	}


}

func TestModifyBmLoadBalancerAttributes(t *testing.T)  {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	params := make(map[string]interface{})
	params["loadBalancerId"] = "lb-3u4w1uir"
	params["loadBalancerName"] = "mudun-test-lb"

	fmt.Println(bmlbService.Call(ModifyBmLoadBalancerAttributes, params))
}

func TestCreateBmListeners(t *testing.T)  {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	params := make(map[string]interface{})
	params["loadBalancerId"] = "lb-5q7g10it"
	params["listeners.0.loadBalancerPort"] = 9876
	params["listeners.0.protocol"] = "tcp"
	params["listeners.0.listenerName"] = "rocketmq"
	params["listeners.0.toaFlag"] = 1
	params["listeners.0.healthSwitch"] = 1

	//{"code":0,"message":"","codeDesc":"Success","requestId":2913648}
	fmt.Println(bmlbService.Call(CreateBmListeners, params))
}

//See: https://cloud.tencent.com/document/product/386/9298
func TestDescribeBmListenerInfo(t *testing.T)  {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	params := make(map[string]interface{})
	params["loadBalancerId"] = "lb-2jrrk40b"

	fmt.Println(bmlbService.Call(DescribeBmListenerInfo, params))


	//fmt.Println(bmlbService.GetL4Listeners("lb-9jd4zn5t"))
}

//See https://cloud.tencent.com/document/product/386/9294
func TestBindBmL4ListenerRs(t *testing.T)  {
	/*
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	bmService := &BmService{
		Service: NewService(Bm, BJ),
	}

	instanceIds := bmService.GetDevices([]string{"10.201.95.2"})

	fmt.Printf("instanceId: %v\n", instanceIds)

	params := make(map[string]interface{})
	params["loadBalancerId"] = "lb-qo7wdyjv"
	params["listenerId"] = "lbl-h0n4a4cj"
	params["backends.0.port"] = 30004
	params["backends.0.instanceId"] = instanceIds[0]
	params["backends.0.weight"] = 10

	fmt.Println(bmlbService.Call(BindBmL4ListenerRs, params))
	*/
}

func TestDescribeBmLoadBalancersTaskResult(t *testing.T)  {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	params := make(map[string]interface{})
	params["requestId"] = 2915276

	resp := bmlbService.Call(DescribeBmLoadBalancersTaskResult, params)
	fmt.Println(resp)

	_, data := CheckSuccess(resp)
	status := data.(map[string]interface{})["status"].(int64)

	fmt.Println(status)

}

func TestGetLbVip(t *testing.T)  {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	fmt.Println(bmlbService.GetLbVip("lb-3u4w1uir"))
}

// See: https://cloud.tencent.com/document/product/386/9304
func TestDeleteBmLoadBalancers(t *testing.T)  {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	params := make(map[string]interface{})
	params["loadBalancerId"] = "lb-2jrrk40b"

	fmt.Println(bmlbService.Call(DeleteBmLoadBalancers, params))
}

func ClearMap(paramMap map[string]interface{}) {
	for key, _ := range paramMap {
		delete(paramMap, key)
	}
}

func TestDescribeLB(t *testing.T) {
	bmlbService := &BmlbService{
		Service: NewService(Bmlb, BJ),
	}

	params := make(map[string]interface{})
	params["loadBalancerIds.0"] = "lb-2jrrk40b"
	params["loadBalancerType"] = "internal"

	fmt.Println(bmlbService.Call(DescribeBmLoadBalancers, params))

}