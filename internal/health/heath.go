package health

import (
	"github.com/labstack/echo/v4"
)

func HealthCheck(c echo.Context) error {
	// TODO: Implement health check logic here
	// Return OK if everything is fine
	// Otherwise, return an error
	return c.JSON(200, "OK")
}
