package di

import "sync"

type any interface{}

type Container struct {
	defines   sync.Map
	instances sync.Map
}

type CreateFunc func(c *Container) interface{}
