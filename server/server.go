package server

import (
	"github.com/gin-gonic/gin"
)


type QcloudServer struct {
	bmRouter *BmRouter
	bmlbRouter *BmlbRouter
	cnsRouter *CnsRouter
	k8sRouter *K8sRouter

	engine *gin.Engine
}

func NewServer() *QcloudServer  {
	service := &QcloudServer{
		engine:gin.Default(),
		bmRouter:NewBmRouter(),
		bmlbRouter:NewBmlbRouter(),
		cnsRouter:NewCnsRouter(),
		k8sRouter: NewK8sRouter(),
	}

	service.Init()

	return service
}

func (server *QcloudServer) Run()()  {
	server.engine.Run(":8360")
}

func (server *QcloudServer) Init() {
	server.engine.GET(server.bmRouter.path, server.bmRouter.handle)
	server.engine.GET(server.bmlbRouter.path, server.bmlbRouter.handle)
	server.engine.GET(server.cnsRouter.path, server.cnsRouter.handle)
	server.engine.GET(server.k8sRouter.path, server.k8sRouter.handle)
}
