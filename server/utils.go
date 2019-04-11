package server

import (
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/json"
	"fmt"
)

type Result struct {
	data string
}

func Resp(msg interface{}) string {
	fmt.Println(msg)
	resp,_ := json.Marshal(msg)
	fmt.Println(resp)
	return fmt.Sprintf("%s", string(resp))
}