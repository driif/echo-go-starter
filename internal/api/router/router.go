package router

import (
	"github.com/driif/echo-go-starter/internal/server"
	mdwr "github.com/driif/echo-go-starter/internal/server/net/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Attaches a router with configurirable middleware and all the routes to the server
func InitGroups(s *server.Server) {
	s.Router = &server.Router{
		// All Available Routes
		Routes: nil,

		Root: s.Echo.Group(""),

		Management: s.Echo.Group("/-", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "query:mgmt-secret",
			Validator: func(key string, c echo.Context) (bool, error) {
				return key == s.Config.Management.Secret, nil
			},
			Skipper: func(c echo.Context) bool {
				switch c.Path() {
				case "/-/ready":
					return true
				}
				return false
			},
		}), mdwr.NoCache()),
	}
}

func AttachRoutes(s *server.Server) {
	// Attach all the routes
	s.Router.Routes = []*echo.Route{
		// == MANAGEMENT == //
		// management.GetVersionRoute(s),
		// management.GetDbVersionRoute(s),
		// == USER == //
		// user.GetMeRoute(s),
		// user.CreateUserRoute(s),
		// user.UpdateUserRoute(s),
		// user.DeleteUserRoute(s),

	}
}
