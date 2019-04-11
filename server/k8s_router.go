package server

import (
	"github.com/gin-gonic/gin"
	"mobike.io/infra/qcloud-agent/qcloud"
	"strings"
	"net/http"
	"fmt"
	"time"
)

const (
	Delete Action = "Delete"
	Create Action = "Create"
)

type K8sRouter struct {
	path string
	handle gin.HandlerFunc
	bmService *qcloud.BmService
	bmlbService *qcloud.BmlbService
	cnsService *qcloud.CnsService
}

func NewK8sRouter() *K8sRouter {
	bmServer := &qcloud.BmService{
		Service:qcloud.NewService(qcloud.Bm, qcloud.BJ),
	}

	bmlbServer := &qcloud.BmlbService{
		Service:qcloud.NewService(qcloud.Bmlb, qcloud.BJ),
	}

	cnsServer := &qcloud.CnsService{
		Service:qcloud.NewService(qcloud.Cns, qcloud.BJ),
	}

	router := &K8sRouter{
		bmService: bmServer,
		bmlbService:bmlbServer,
		cnsService:cnsServer,
		path: "/k8s",
	}

	router.init()

	return router
}


/*
http://127.0.0.1:8360/k8s?
	action=Create
	&subDomain=zk.ydt
	&lbName=test-ydt-lb
	&ports.0=9092:30002
	&ports.1=9093:30003
	&nodeIps.0=10.201.95.2
	&nodeIps.1=10.201.53.2

http://127.0.0.1:8360/k8s?
	action=Delete
	&subDomain=zk.ydt
	&lbName=test-ydt-lb
*/
func (router *K8sRouter) init()  {
	router.handle = func(context *gin.Context){
		action := Action(context.Query("action"))
		switch action {
		case Create:
			query := context.Request.URL.Query()
			subDomain := context.Query("subDomain")
			lbName := context.Query("lbName")
			portsMap := []string{}
			nodeIps := []string{}

			// 0. get the ports & nodeIps
			for k, v := range query {
				if strings.Contains(k, "ports.") {
					portsMap = append(portsMap, v[0])
				} else if strings.Contains(k, "nodeIps.") {
					nodeIps = append(nodeIps, v[0])
				}
			}

			if portsMap == nil || nodeIps == nil || lbName == "" || subDomain == "" {
				context.JSON(http.StatusOK, "params invalid")
				return
			}

			params := make(map[string]interface{})

			// 1. get nodeIds for the nodeIps
			for idx, value := range nodeIps {
				k := fmt.Sprintf("lanIps.%v", idx)
				params[k] = value
			}
			nodeIds := router.bmService.GetDevices(params)
			fmt.Printf("get the instances: %v\n", nodeIds)

			// 2. get lb if exists
			loadBalancerId, err := router.bmlbService.GetLB("internal", lbName)
			fmt.Printf("lbName: %s\n", lbName)
			if loadBalancerId == "" || err != nil {
				fmt.Printf("lb %v is no exists, let's create it\n", lbName)

				params["nodeIp"] = nodeIps[0]
				params["loadBalancerType"] = "internal"
				params["goodsNum"] = 1

				loadBalancerId, _ = router.bmlbService.CreateLB(params)
				time.Sleep(time.Second * 10)
				router.bmlbService.ModifyLbName(loadBalancerId, lbName)
				fmt.Printf("successed to created lb: %v\n", loadBalancerId)
			}

			time.Sleep(time.Second * 10)
			// 3. create lb listener
			ClearMap(params)
			params["loadBalancerId"] = loadBalancerId
			for i, v := range portsMap {
				//9092:30001,9093:30002
				portMap := strings.Split(v, ":")
				pre := fmt.Sprintf("listeners.%v.", i)
				params[pre+"loadBalancerPort"] = portMap[0]
				params[pre+"protocol"] = "tcp"
				params[pre+"listenerName"] = portMap[0]
				params[pre+"toaFlag"] = 1
				params[pre+"healthSwitch"] = 1
			}
			router.bmlbService.CreateBmListeners(params)

			// 4. get listeners and bind rs
			listeners, _ := router.bmlbService.GetL4Listeners(loadBalancerId)
			fmt.Printf("successed to create lb listeners: %v\n", listeners)
			for _, v := range portsMap {
				//9092:30001,9093:30002
				portMap := strings.Split(v, ":")
				ClearMap(params)
				params["loadBalancerId"] = loadBalancerId
				params["listenerId"] = listeners[portMap[0]]
				index := 0
				for _, instanceId := range nodeIds {
					idx := fmt.Sprintf("backends.%v.", index)
					params[idx+"port"] = portMap[1]
					params[idx+"instanceId"] = instanceId
					params[idx+"weight"] = 10
				}
				router.bmlbService.Call(qcloud.BindBmL4ListenerRs, params)
				time.Sleep(time.Second * 5)
				fmt.Printf("successed to bind rs: %v", params)
			}

			// 5. create record
			lbVip, _ := router.bmlbService.GetLbVip(loadBalancerId)
			fmt.Printf("lbvip : %v\n", lbVip)
			ClearMap(params)
			params["domain"] = "bike.io"
			params["subDomain"] = subDomain
			params["recordType"] = "A"
			params["recordLine"] = "默认"
			params["value"] = lbVip
			params["ttl"] = 10
			router.cnsService.Call(qcloud.RecordCreate, params)
			fmt.Printf("success to create record", params)

			context.JSON(http.StatusOK, 0)

		case Delete:
			subDomain := context.Query("subDomain")
			lbName := context.Query("lbName")
			params := make(map[string]interface{})

			// 1. get and delete lb
			loadBalancerId, _ := router.bmlbService.GetLB("internal", lbName)
			if loadBalancerId != ""  {
				params["loadBalancerId"] = loadBalancerId
				router.bmlbService.Call(qcloud.DeleteBmLoadBalancers, params)
			}

			// 2. get and delete record
			recordId, _ := router.cnsService.GetRecord("bike.io", subDomain)
			fmt.Printf("%v record id : %v\n", subDomain, recordId)
			if recordId != -1 {
				ClearMap(params)
				params["domain"] = "bike.io"
				params["recordId"] = recordId

				router.cnsService.Call(qcloud.RecordDelete, params)
			}

			context.JSON(http.StatusOK, 0)
		}
	}
}

func ClearMap(paramMap map[string]interface{}) {
	for key, _ := range paramMap {
		delete(paramMap, key)
	}
}
