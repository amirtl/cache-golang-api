package main

import (
	"training/internal/pkg/service"
)

func main() {
	arg := service.LoadConfig()
	//service.Engine(arg).Run(arg.Conf.BaseUrl)
	service.Engine(arg).Run(arg.Conf.BaseUrl)
}
