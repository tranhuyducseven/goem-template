package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func RequestLogMiddleware() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogLatency: true,
		LogMethod:  true,
		LogError:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Method == "OPTIONS" {
				return nil
			}

			var e *zerolog.Event
			if v.Status >= 400 {
				e = log.Err(v.Error)
			} else {
				e = log.Info()
			}

			e.Int("status", v.Status).
				Str("method", v.Method).
				Str("uri", v.URI).
				Str("latency", v.Latency.String()).
				Send()
			return nil
		},
	})
}
