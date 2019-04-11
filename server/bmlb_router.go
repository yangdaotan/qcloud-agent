package server

import (
	"github.com/gin-gonic/gin"
	"mobike.io/infra/qcloud-agent/qcloud"
	"net/http"
	"time"
	"strings"
)

const(
	//bmlb
	CreateBmLoadBalancer Action = "CreateBmLoadBalancer"
	ModifyBmLoadBalancerAttributes Action = "ModifyBmLoadBalancerAttributes"
	CreateBmListeners Action = "CreateBmListeners"
	BindBmL4ListenerRs Action = "BindBmL4ListenerRs"
)


type BmlbRouter struct {
	path string
	handle gin.HandlerFunc
	service *qcloud.BmlbService
}

func NewBmlbRouter() *BmlbRouter {
	service := &qcloud.BmlbService{
		Service:qcloud.NewService(qcloud.Bmlb, qcloud.BJ),
	}

	router := &BmlbRouter{
		service: service,
		path: "/bmlb",
	}

	router.init()

	return router
}

func (router *BmlbRouter) init()  {
	router.handle = func(context *gin.Context){
		action := Action(context.Query("action"))
		var resp interface{}
		switch action {
		case CreateBmLoadBalancer:
			nodeIp := context.Query("nodeIp")
			loadBalancerType := context.DefaultQuery("loadBalancerType", "internal")
			goodsNum := context.DefaultQuery("goodsNum", "1")
			loadBalancerName := context.Query("loadBalancerName")
			params := make(map[string]interface{})
			params["nodeIp"] = nodeIp
			params["loadBalancerType"] =loadBalancerType
			params["goodsNum"] = goodsNum
			resp, _= router.service.CreateLB(params)
			if loadBalancerName != "" {
				time.Sleep(time.Second * 10)
				router.service.ModifyLbName(resp.(string), loadBalancerName)
			}
		case ModifyBmLoadBalancerAttributes:
			loadBalancerId := context.Query("loadBalancerId")
			loadBalancerName := context.Query("loadBalancerName")
			resp, _ = router.service.ModifyLbName(loadBalancerId, loadBalancerName)
		case CreateBmListeners:
			query := context.Request.URL.Query()
			loadBalancerId := context.Query("loadBalancerId")
			params := make(map[string]interface{})
			params["loadBalancerId"] = loadBalancerId
			for k, v := range query {
				if strings.Contains(k, "listeners.") {
					params[k] = v[0]
				}
			}
			resp, _ = router.service.CreateBmListeners(params)
		case BindBmL4ListenerRs:
			query := context.Request.URL.Query()
			params := make(map[string]interface{})
			params["loadBalancerId"] = context.Query("loadBalancerId")
			params["listenerId"] = context.Query("listenerId")
			for k, v := range query {
				if strings.Contains(k, "backends.") {
					params[k] = v[0]
				}
			}
			context.String(http.StatusOK, router.service.Call(qcloud.BindBmL4ListenerRs, params))
			return
		}

		context.JSON(http.StatusOK, resp)
	}
}
