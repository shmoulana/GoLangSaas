package main

import (
	"github.com/shmoulana/Redios/cmd/webservice"
	"github.com/shmoulana/Redios/configs"
)

func main() {
	// initialize config
	configs.Init()
	conf := configs.Get()

	webservice.StartServer(*conf)
}
