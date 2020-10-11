# echoboot

boot 是一个能够简化 Echo 框架下路由管理的库

核心思想是通过 boot 注册 Echo 的实例和路由群组，然后在每个 Controller 的 init 中各自注册自己的路由回调函数，再统一由 boot 控制群组和路由回调函数的加载顺序，减少路由间的耦合

boot 是非线程安全的（因为 boot 的业务需要在 Echo 服务启动前加载并运行完毕，所以无需考虑并发，也不建议在此期间出现并发问题）

# Installation

```
$ go get -u github.com/gokit/echo
```

# Quick start

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/gokit/echox/boot"
    _ "github.com/gokit/echox/boot/example/controller" // 加载 controller 包下的 init 函数，注册路由相关的回调函数
)

func main() {
    e := echo.New()

    // 注册默认的 Fiber 实例并运行
    boot.Strap(e)

    e.Logger.Fatal(e.Start(":8080"))
}
```

controller

```go
package controller

import (
    "github.com/gokit/echox/boot"
    "github.com/labstack/echo/v4"
)

type IndexController struct{}

func init() {
    // 向默认的 Gin 中注册路由
    boot.Route(new(IndexController).registerRoutes)
}

func (i *IndexController) registerRoutes(app *echo.Echo) {
    app.GET("/", i.index)
}

func (IndexController) index(c echo.Context) error {
    return c.String(200, "hello world")
}
```

注册路由群组

```go
package controller

import (
    "github.com/gokit/echox/boot"
    "github.com/labstack/echo/v4"
)

const GroupUser = "/user"
const GroupAdmin = "/admin"

func init() {
    boot.Group(func(app *echo.Echo) {
        boot.AddGroup(GroupUser, app.Group("/user"))
        boot.AddGroup(GroupAdmin, app.Group("/admin"))
    })
}
```

更多其他相关方法

```go
// 注册一个指定 name 的 Group 实例
func AddGroup(name string, group *echo.Group)

// 向指定的 Echo 实例中注入中间件回调
func Middleware(callback func(app *echo.Echo))

// 向指定的 Echo 实例中注入路由回调
func Route(callback func(app *echo.Echo))

// 向指定的 Echo 实例中注入路由群组回调
func Group(callback func(app *echo.Echo))

// 向指定 name 的路由群组中中注入路由回调
func GroupByName(groupName string, callback func(group *echo.Group))

// 在注册 Echo 实例后，进行初始化，boot 会优先调用路由群组的回调函数，然后再调用路由的回调函数
func Init(app *echo.Echo)
```
