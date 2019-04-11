package main

import (
	"mobike.io/infra/qcloud-agent/server"
)


func main()  {
	server := server.NewServer()
	server.Run()
}