package qcloud

import (
	"testing"
	"fmt"
)

func TestSig(t *testing.T) {

	params := make(map[string]interface{})

	params["Action"] = "DescribeInstances"
	params["Nonce"] = 11886
	params["Region"] = "gz"
	params["SecretId"] = "AKIDz8krbsJ5yKBZQpn74WFkmLPx3gnPhESA"
	params["Timestamp"] = 1465185768
	params["eipId"] = "eip-testcpm"
	params["instanceId"] = "cpm-09dx96dg"

	signature := GenerateSignature("GET", "bmeip", "/v2/index.php", params)

	fmt.Println(signature)
}

func TestCurl(t *testing.T)  {

	resp := Get("http://baidu.com", nil)

	fmt.Println(resp)
}


func TestDescribeBmLoadBalancers(t *testing.T)  {

	params := GetCommonParamsMap("DescribeBmLoadBalancers", "bj")

	params["loadBalancerType"] = "internal"
	params["offset"] = 0
	params["limit"] = 10

	signature := GenerateSignature("GET", "bmlb", "/v2/index.php", params)
	fmt.Println(signature)

	params["Signature"] = signature
	domain := "https://bmlb.api.qcloud.com/v2/index.php"

	resp := Get(domain, params)

	fmt.Println(resp)
}

