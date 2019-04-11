package server

import (
	"github.com/gin-gonic/gin"
	"mobike.io/infra/qcloud-agent/qcloud"
	"net/http"
	"strings"
)

const(
	//lm
	GetDeviceInstanceIds Action = "GetDeviceInstanceIds"
)


type BmRouter struct {
	path string
	handle gin.HandlerFunc
	service *qcloud.BmService
}

func NewBmRouter() *BmRouter {
	bmServer := &qcloud.BmService{
		Service:qcloud.NewService(qcloud.Bm, qcloud.BJ),
	}

	router := &BmRouter{
		service: bmServer,
		path: "/bm",
	}

	router.init()

	return router
}

func (router *BmRouter) init()  {
	router.handle = func(context *gin.Context){
		action := Action(context.Query("action"))
		var resp interface{}
		switch action {
		case GetDeviceInstanceIds:
			query := context.Request.URL.Query()
			lanIps := make(map[string]interface{})
			for k, v := range query {
				if strings.Contains(k, "lanIps.") {
					lanIps[k] = v[0]
				}
			}
			resp = router.service.GetDevices(lanIps)
		}

		context.JSON(http.StatusOK, resp)
	}
}
