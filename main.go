package main

import (
	"log"
	"os"
	"./di"
	"./GSR"
	"fmt"
)

type C struct {
	Logger3 GSR.Logger `inject:"auto"`
}

type B struct {
	My  *C     `inject:"auto"`
	Env string `inject:"env"`
}

func main() {
	container := di.NewContainer().Register("env", "test").Register(
		(*GSR.Logger)(nil),
		di.Create(func(c *di.Container) interface{} {
			return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
		}),
	).Register(
		(*C)(nil),
		di.Create(func(c *di.Container) interface{} {
			return &C{}
		}),
	).Register(
		(*B)(nil),
		di.Create(func(c *di.Container) interface{} {
			return &B{}
		}),
	)

	//b := B{}
	//container.InjectOn(&b)
	b := container.Get((*B)(nil)).(*B)
	b.My.Logger3.Println("haha")
	fmt.Println(b.Env)
}
