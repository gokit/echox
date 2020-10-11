package sessions

import (
	ctx "context"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

const (
	DefaultKey  = "github.com/gokit/echox/sessions"
	errorFormat = "[sessions] ERROR! %s\n"
)

type Store interface {
	sessions.Store
	Options(Options)
}

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Session interface {
	// Get Session ID
	ID() string
	// Get returns the session value associated to the given key.
	Get(key interface{}) interface{}
	// Set sets the session value associated to the given key.
	Set(key interface{}, val interface{})
	// Delete removes the session value associated to the given key.
	Delete(key interface{})
	// Clear deletes all values in the session.
	Clear()
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string)
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}
	// Options sets configuration for a session.
	Options(Options)
	// Save saves all sessions used during the current request.
	Save() error
	// Destory session
	Destory() error
	// Session is new
	IsNew() bool
}

// gin
//func Sessions(name string, store Store) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		s := &session{name, c.Request, store, nil, false, c.Writer}
//		c.Set(DefaultKey, s)
//		defer context.Clear(c.Request)
//		c.Next()
//		// Save except cookie store
//		s.Save()
//	}
//}

// echo
func Sessions(name string, store Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			s := &session{name, c.Request(), store, nil, false, c.Response()}
			c.Set(DefaultKey, s)
			defer context.Clear(c.Request())
			err := next(c)

			if err != nil {
				return err
			}

			err = s.Save()

			if err != nil {
				c.Logger().Errorf(errorFormat, err)
			}

			return err
		}
	}
}

// echo - store session by context
func SessionsByContext(name string, store Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			s := &session{name, c.Request(), store, nil, false, c.Response()}
			c.SetRequest(c.Request().WithContext(ctx.WithValue(c.Request().Context(), DefaultKey, s)))
			defer context.Clear(c.Request())
			err := next(c)

			if err != nil {
				return err
			}

			err = s.Save()

			if err != nil {
				c.Logger().Errorf(errorFormat, err)
			}

			return err
		}
	}
}

// gin
//func SessionsMany(names []string, store Store) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		sessions := make(map[string]Session, len(names))
//		for _, name := range names {
//			sessions[name] = &session{name, c.Request, store, nil, false, c.Writer}
//		}
//		c.Set(DefaultKey, sessions)
//		defer context.Clear(c.Request)
//		c.Next()
//		// Save except cookie store
//		for _, name := range names {
//			sessions[name].Save()
//		}
//	}
//}

// echo
func SessionsMany(names []string, store Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessions := make(map[string]Session, len(names))
			for _, name := range names {
				sessions[name] = &session{name, c.Request(), store, nil, false, c.Response()}
			}
			c.Set(DefaultKey, sessions)
			defer context.Clear(c.Request())
			err := next(c)

			if err != nil {
				return err
			}

			// Save except cookie store
			for _, name := range names {
				err = sessions[name].Save()
				if err != nil {
					break
				}
			}

			if err != nil {
				c.Logger().Errorf(errorFormat, err)
			}

			return err
		}
	}
}

// echo - store session by context
func SessionsManyByContext(names []string, store Store) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sessions := make(map[string]Session, len(names))
			for _, name := range names {
				sessions[name] = &session{name, c.Request(), store, nil, false, c.Response()}
			}
			c.SetRequest(c.Request().WithContext(ctx.WithValue(c.Request().Context(), DefaultKey, sessions)))
			defer context.Clear(c.Request())
			err := next(c)

			if err != nil {
				return err
			}

			// Save except cookie store
			for _, name := range names {
				err = sessions[name].Save()
				if err != nil {
					break
				}
			}

			if err != nil {
				c.Logger().Errorf(errorFormat, err)
			}

			return err
		}
	}
}

type session struct {
	name    string
	request *http.Request
	store   Store
	session *sessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *session) ID() string {
	return s.Session().ID
}

func (s *session) Get(key interface{}) interface{} {
	return s.Session().Values[key]
}

func (s *session) Set(key interface{}, val interface{}) {
	s.Session().Values[key] = val
	s.written = true
}

func (s *session) Delete(key interface{}) {
	delete(s.Session().Values, key)
	s.written = true
}

func (s *session) Clear() {
	for key := range s.Session().Values {
		s.Delete(key)
	}
}

func (s *session) AddFlash(value interface{}, vars ...string) {
	s.Session().AddFlash(value, vars...)
	s.written = true
}

func (s *session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *session) Options(options Options) {
	s.Session().Options = options.ToGorillaOptions()
}

func (s *session) Save() error {
	if s.Written() || s.Session().IsNew {
		e := s.Session().Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

func (s *session) Destory() error {
	s.Options(Options{MaxAge: -1})
	s.written = true
	return s.Save()
}

func (s *session) IsNew() bool {
	return s.Session().IsNew
}

func (s *session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		} else if s.session.IsNew {
			// Save if session is new before write response
			s.Save()
		}
	}
	return s.session
}

func (s *session) Written() bool {
	return s.written
}

// shortcut to get session
//func Default(c *gin.Context) Session {
//	return c.MustGet(DefaultKey).(Session)
//}

// shortcut to get session
func Default(c echo.Context) Session {
	sess := c.Get(DefaultKey)
	if sess != nil {
		return sess.(Session)
	}
	return c.Request().Context().Value(DefaultKey).(Session)
}

// shortcut to get session from context
func DefaultByContext(r *http.Request) Session {
	return r.Context().Value(DefaultKey).(Session)
}

// shortcut to get session with given name
//func DefaultMany(c *gin.Context, name string) Session {
//	return c.MustGet(DefaultKey).(map[string]Session)[name]
//}

// shortcut to get session with given name
func DefaultMany(c echo.Context, name string) Session {
	sess := c.Get(DefaultKey)
	if sess != nil {
		return sess.(map[string]Session)[name]
	}
	return c.Request().Context().Value(DefaultKey).(map[string]Session)[name]
}

// shortcut to get session with given name from context
func DefaultManyByContext(r *http.Request, name string) Session {
	return r.Context().Value(DefaultKey).(map[string]Session)[name]
}
