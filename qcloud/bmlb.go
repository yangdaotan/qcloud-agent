package qcloud

import (
	"errors"
	"fmt"
	"time"
)


type BmlbService struct {
	Service *Service
}

// See: https://cloud.tencent.com/document/product/386/9303
func (bmlbService *BmlbService) CreateLB(params map[string]interface{}) (string, error) {

	bmvpcService := &BmvpcService{
		Service: NewService(Bmvpc, BJ),
	}

	nodeIp := params["nodeIp"].(string)

	vpc := bmvpcService.getMatchVpc(nodeIp)

	if vpc == nil || vpc.subnetNum < 1 {
		return "", errors.New("no avaliable vpc to create lb")
	}

	unSubnetId := bmvpcService.getAvailableSubnet(vpc.unVpcId)

	if unSubnetId == "" {
		return "", errors.New("no avaliable subnet to create lb")
	}

	params["unSubnetId"] = unSubnetId
	params["unVpcId"] = vpc.unVpcId
	delete(params, "nodeIp")

	resp := bmlbService.Service.call(CreateBmLoadBalancer, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return "", errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	loadBalancerId := dat["loadBalancerIds"].([]interface{})[0].(string)

	return  loadBalancerId, nil
}

func (bmlbService *BmlbService) CreateBmListeners(params map[string]interface{}) (int64, error) {
	resp := bmlbService.Service.call(CreateBmListeners, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	fmt.Println(dat)
	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return -1, errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	requestId := dat["requestId"].(int64)

	status := bmlbService.await(requestId)

	if status == 1 || status == -1 {
		return status, errors.New("failed to create mbListeners")
	}

	return status, nil
}

func (bmlbService *BmlbService) CreateBmForwardListeners(params map[string]interface{}) (int64, error) {
	resp := bmlbService.Service.call(CreateBmForwardListeners, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	fmt.Println(dat)
	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return -1, errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	requestId := dat["requestId"].(int64)

	status := bmlbService.await(requestId)

	if status == 1 || status == -1 {
		return status, errors.New("failed to create mbListeners")
	}

	return status, nil
}

func (bmlbService *BmlbService) GetL4Listeners(loadBalancerId string) (map[string]string, error) {
	params := make(map[string]interface{})
	params["loadBalancerId"] = loadBalancerId

	resp := bmlbService.Call(DescribeBmListenerInfo, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return nil, errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	listeners := make(map[string]string)
	for _, listener := range dat["listenerSet"].([]interface{}) {
		listenerName := listener.(map[string]interface{})["listenerName"].(string)
		listeners[listenerName] = listener.(map[string]interface{})["listenerId"].(string)
	}

	return listeners, nil
}

func (bmlbService *BmlbService) GetL7ListenerId(loadBalancerId string) (string, error) {
	params := make(map[string]interface{})
	params["loadBalancerId"] = loadBalancerId

	resp := bmlbService.Call(DescribeBmForwardListenerInfo, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return "", errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	return dat["listenerSet"].([]interface{})[0].(map[string]interface{})["listenerId"].(string), nil
}

func (bmlbService *BmlbService) GetLbVip(loadBalanceId string) (string, error) {
	params := make(map[string]interface{})
	params["loadBalancerType"] = "internal"
	params["loadBalancerIds.0"] = loadBalanceId

	resp := bmlbService.Call(DescribeBmLoadBalancers, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return "", errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	return dat["loadBalancerSet"].([]interface{})[0].(map[string]interface{})["loadBalancerVips"].([]interface{})[0].(string), nil
}

func (bmlbService *BmlbService) ModifyLbName(loadbalancerId, name string) (int64, error) {
	params := make(map[string]interface{})
	params["loadBalancerId"] = loadbalancerId
	params["loadBalancerName"] = name

	resp := bmlbService.Call(ModifyBmLoadBalancerAttributes, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("resp: %v\n", dat)
	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return -1, errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	requestId := dat["requestId"].(int64)

	status := bmlbService.await(requestId)

	if status == 1 || status == -1 {
		return status, errors.New("failed to create mbListeners")
	}

	return status, nil
}

func (bmlbService *BmlbService) Call(action Action, params map[string]interface{}) string {
	return bmlbService.Service.call(action, params)
}

func (bmlbService *BmlbService) await(requestId int64) int64 {

	for {
		params := make(map[string]interface{})
		params["requestId"] = requestId

		resp := bmlbService.Call(DescribeBmLoadBalancersTaskResult, params)
		fmt.Printf("status: %v\n", resp)

		err, data := CheckSuccess(resp)
		if err != nil {
			return -1
		}

		status := data.(map[string]interface{})["status"].(int64)

		if status == 2 {
			time.Sleep(time.Second * 2)
		} else {
			return status
		}
	}
}


func (bmlbService *BmlbService) GetLB(loadBalancerType, name string) (string, error) {
	params := make(map[string]interface{})
	params["loadBalancerType"] = loadBalancerType
	params["loadBalancerName"] = name

	resp := bmlbService.Call(DescribeBmLoadBalancers, params)

	var dat map[string]interface{}
	if err := Unmarshal([]byte(resp), &dat); err != nil {
		fmt.Println(err)
	}

	fmt.Printf("resp: %v\n", dat)
	code, _ := dat["code"].(int)
	if code != SUCCESS {
		return "", errors.New(fmt.Sprintf("%v:%v", code, dat["message"]))
	}

	loadBalancerSet := dat["loadBalancerSet"].([]interface{})

	if loadBalancerSet == nil || len(loadBalancerSet) == 0 {
		return "", nil
	}

	return loadBalancerSet[0].(map[string]interface{})["loadBalancerId"].(string), nil
}