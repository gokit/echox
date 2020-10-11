package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/filesystem"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	store := filesystem.NewStore("./", []byte("secret"))
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

		session.Destory()

		c.JSON(200, echo.Map{"count": count})
		return nil
	})
	r.Start(":8000")
}
