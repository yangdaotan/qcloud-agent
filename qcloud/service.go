package qcloud

import "fmt"

type Service struct {
	Module Module
	Region Region
}

func NewService(module Module, region Region) *Service {
	return &Service{
		Module:module,
		Region:region,
	}
}

func (service* Service) call(action Action, params map[string]interface{}) string {
	curlParams := GetCommonParamsMap(action, service.Region)

	if params != nil {
		for key, value := range params {
			curlParams[key] = value
		}
	}

	signature := GenerateSignature(Method, service.Module, HostPath, curlParams)
	curlParams["Signature"] = signature

	//"https://bmlb.api.qcloud.com/v2/index.php"
	domain := fmt.Sprintf("https://%v.%v%v", service.Module, DomainBase, HostPath)

	return Get(domain, curlParams)
}

