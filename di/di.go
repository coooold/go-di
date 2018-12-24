package di

import (
	"reflect"
	"fmt"
)

func New() *Container {
	return &Container{}
}

// 获得真实的key
func getKey(t any) any {
	r := reflect.ValueOf(t)
	var k interface{}
	switch r.Kind() {
	case reflect.Ptr:
		k = r.Type().Elem()
	case reflect.String:
		k = t
	default:
		panic(fmt.Errorf("t can only be typeOf ptr or interface, %s was given", r.Kind()))
	}

	return k
}

// Register 注册值或构造函数
func (c *Container) Register(t any, def any) *Container {
	k := getKey(t)

	if _, ok := def.(CreateFunc); ok { // 如果是构造函数，那么放到defines
		c.defines.Store(k, def)
	} else { // 其他的直接存值
		c.instances.Store(k, def)
	}
	return c
}

// Create 用来包装构造函数
func Create(cf func(c *Container) interface{}) CreateFunc {
	return cf
}

// 获取实例
func (c *Container) Get(t any) interface{} {
	var k any

	switch t.(type) {
	case reflect.Type:
		k = t
	case string:
		k = t
	default:
		k = getKey(t)
	}

	if v, ok := c.instances.Load(k); ok {
		return v
	}

	if def, ok := c.defines.Load(k); ok {
		if cf, ok := def.(CreateFunc); ok { // 如果是构造函数
			v := cf(c)
			c.InjectOn(v)
			c.instances.Store(k, v)
			return v
		}
	}

	return nil
}

// Inject 注入，只能注入到指针
func (c *Container) InjectOn(p interface{}) {
	v := reflect.ValueOf(p)

	if v.Kind() != reflect.Ptr {
		panic("param should be ptr")
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("inject")

		fieldValue := v.FieldByName(field.Name)

		if !fieldValue.CanSet() || len(tag) == 0 {
			continue
		}

		var k any
		switch field.Type.Kind() {
		case reflect.Ptr:
			k = field.Type.Elem()
		case reflect.Interface:
			k = field.Type
		default:
			k = tag // 注入值
		}

		instance := c.Get(k)
		if instance == nil {
			continue
		}
		fieldValue.Set(reflect.ValueOf(instance))
	}
}
