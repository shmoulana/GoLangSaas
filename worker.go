package main

import (
	"fmt"
)

type A interface {
	Name() string
	Type() string
}

type B struct {
}

func (t B) Name() string {
	return "test"
}
func (t B) Type() string {
	return "test"
}

func main() {
	var a interface{}
	a = B{}

	fmt.Print(a.(A).Name())
	// configs.Init()
	// conf := configs.Get()
	// fac := internal.Transport{}

	// queueService := fac.GetQueueService(*conf)

	// queueService.Run(context.Background())
}
