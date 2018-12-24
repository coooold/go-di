# go-di
The dependency injection container for golang with auto injection

## Auto injection
```Go
type A struct {
	Logger GSR.Logger `inject:"auto"`
}

func (a *A) writeLog(msg string) {
    a.Logger.Println(msg)
}
```

## Make a new container
```Go
container := di.New()
```

## Configure the container

```Go
// Register a value
container.Register("env", "product")

// Register an interface
container.Register(
    (*GSR.Logger)(nil),
    di.Create(func(c *di.Container) interface{} {
        return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
    }),
)

// Register a struct
container.Register(
    (*A)(nil),
    di.Create(func(c *di.Container) interface{} {
        return &A{
            StartTime: c.Get("StartTime"),
        }
    }),
)

// Register an instance
container.Register(
    (*B)(nil),
    &B{},
)
```

## Fetch the instance or value
```Go
container.Get((*A)(nil))
container.Get("env")

```

## Call in chain
```Go
container.Register("a", 1).Register("b", 2).Register("c", 3)

```

## Overview

```Go

type C struct {
	Logger3 GSR.Logger `inject:"auto"`
}

type B struct {
	My  *C     `inject:"auto"`
	Env string `inject:"env"`
}

func main() {
	container := di.New().Register("env", "test").Register(
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



```