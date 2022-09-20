package main

import (
	"github.com/shmoulana/Redios/cmd/webservice"
	"github.com/shmoulana/Redios/configs"
)

func main() {
	// initialize config
	conf := configs.Init()

	webservice.StartServer(*conf)
}
