package bind

import (
	"github.com/driif/echo-go-starter/internal/server/net/runtime"
	"github.com/labstack/echo/v4"
)

// Body binds the request body to the given Validatable and validates it.
func Body(c echo.Context, v runtime.Validatable) error {
	binder := c.Echo().Binder.(*echo.DefaultBinder)

	if err := binder.BindBody(c, v); err != nil {
		return err
	}

	return v.Validate()
}
