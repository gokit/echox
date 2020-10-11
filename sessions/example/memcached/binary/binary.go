package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/memcached"
	"github.com/labstack/echo/v4"
	"github.com/memcachier/mc"
)

func main() {
	r := echo.New()
	client := mc.NewMC("localhost:11211", "username", "password")
	store := memcached.NewMemcacheStore(client, "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/incr", func(c echo.Context) error {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, echo.Map{"count": count})
		return nil
	})
	r.Start(":8000")
}
