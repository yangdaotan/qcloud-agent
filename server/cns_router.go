package server

import (
	"github.com/gin-gonic/gin"
	"mobike.io/infra/qcloud-agent/qcloud"
	"net/http"
)

const(
	//cns
	RecordCreate Action = "RecordCreate"
)


type CnsRouter struct {
	path string
	handle gin.HandlerFunc
	service *qcloud.CnsService
}

func NewCnsRouter() *CnsRouter {
	cnsServer := &qcloud.CnsService{
		Service:qcloud.NewService(qcloud.Cns, qcloud.BJ),
	}

	router := &CnsRouter{
		service: cnsServer,
		path: "/cns",
	}

	router.init()

	return router
}

func (router *CnsRouter) init()  {
	router.handle = func(context *gin.Context){
		action := Action(context.Query("action"))
		var resp interface{}
		switch action {
		case RecordCreate:
			params := make(map[string]interface{})
			params["domain"] = "bike.io"
			params["subDomain"] = context.Query("subDomain")
			params["recordType"] = "A"
			params["recordLine"] = "默认"
			params["value"] = context.Query("vip")
			params["ttl"] = 10

			resp = router.service.Call(qcloud.RecordCreate, params)
		}

		context.JSON(http.StatusOK, resp)
	}
}
