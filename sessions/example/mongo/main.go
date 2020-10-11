package main

import (
	"github.com/globalsign/mgo"
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/mongo"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	session, err := mgo.Dial("localhost:27017/test")
	if err != nil {
		// handle err
	}

	c := session.DB("").C("sessions")
	store := mongo.NewStore(c, 3600, true, []byte("secret"))
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
	})
	r.Start(":8000")
}
