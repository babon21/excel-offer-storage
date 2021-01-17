package middleware

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/rs/zerolog/log"
	"time"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

// CORS will handle the CORS middleware
func (m *GoMiddleware) AccessLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		log.Info().Msg(fmt.Sprintf("[%s] %s, %s %s\n", c.Request().Method, c.Request().RemoteAddr,
			c.Request().URL.Path, time.Since(start)))
		return err
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
