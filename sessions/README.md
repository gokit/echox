# sessions

[![Go Report Card](https://goreportcard.com/badge/github.com/gokit/echox/sessions)](https://goreportcard.com/report/github.com/gokit/echox/sessions)
[![GoDoc](https://godoc.org/github.com/gokit/echox/sessions?status.svg)](https://godoc.org/github.com/gokit/echox/sessions)

Echo middleware for session management with multi-backend support:

- [cookie-based](#cookie-based)
- [filesystem](#filesystem)
- [Redis](#redis)
- [memcached](#memcached)
- [MongoDB](#mongodb)
- [memstore](#memstore)

## Usage

### Start using it

Download and install it:

```bash
$ go get github.com/gokit/echox/sessions
```

Import it in your code:

```go
import "github.com/gokit/echox/sessions"
```

## Basic Examples

### single session

```go
package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/cookie"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(c echo.Context) error {
		session := sessions.Default(c)

		if session.Get("hello") != "world" {
			session.Set("hello", "world")
			session.Save()
		}

		c.JSON(200, echo.Map{"hello": session.Get("hello")})

        return nil
	})
	r.Start(":8000")
}
```

### multiple sessions

```go
package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/cookie"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	store := cookie.NewStore([]byte("secret"))
	sessionNames := []string{"a", "b"}
	r.Use(sessions.SessionsMany(sessionNames, store))

	r.GET("/hello", func(c echo.Context) error {
		sessionA := sessions.DefaultMany(c, "a")
		sessionB := sessions.DefaultMany(c, "b")

		if sessionA.Get("hello") != "world!" {
			sessionA.Set("hello", "world!")
			sessionA.Save()
		}

		if sessionB.Get("hello") != "world?" {
			sessionB.Set("hello", "world?")
			sessionB.Save()
		}

		c.JSON(200, echo.Map{
			"a": sessionA.Get("hello"),
			"b": sessionB.Get("hello"),
		})
        
        return nil
	})
	r.Start(":8000")
}
```

## Backend examples

### cookie-based

[embedmd]:# (example/cookie/main.go go)
```go
package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/cookie"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	store := cookie.NewStore([]byte("secret"))
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
```


### filesystem

[embedmd]:# (example/filesystem/main.go go)
```go
package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/filesystem"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	store := filesystem.NewStore("/path/to/sessions/", []byte("secret"))
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
```

### Redis

[embedmd]:# (example/redis/main.go go)
```go
package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/redis"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
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
```

### Memcached

#### ASCII Protocol

[embedmd]:# (example/memcached/ascii/ascii.go go)
```go
package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/memcached"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	store := memcached.NewStore(memcache.New("localhost:11211"), "", []byte("secret"))
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
```

#### Binary protocol (with optional SASL authentication)

[embedmd]:# (example/memcached/binary/binary.go go)
```go
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
```

### MongoDB

[embedmd]:# (example/mongo/main.go go)
```go
package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/mongo"
	"github.com/labstack/echo/v4"
	"github.com/globalsign/mgo"
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
        return nil
	})
	r.Start(":8000")
}
```

### memstore

[embedmd]:# (example/memstore/main.go go)
```go
package main

import (
	"github.com/gokit/echox/sessions"
	"github.com/gokit/echox/sessions/memstore"
	"github.com/labstack/echo/v4"
)

func main() {
	r := echo.New()
	store := memstore.NewStore([]byte("secret"))
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
```


