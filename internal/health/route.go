package health

import (
	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, path string) {
	e.GET(path, HealthCheck)
}
