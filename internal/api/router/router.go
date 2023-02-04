package router

// import (
// 	"github.com/cubular-io/go-chat/api/handlers/management"
// 	"github.com/cubular-io/go-chat/api/handlers/msg"
// 	"github.com/cubular-io/go-chat/api/handlers/room"
// 	"github.com/cubular-io/go-chat/api/handlers/user"
// 	"github.com/cubular-io/go-chat/server"
// 	"github.com/cubular-io/go-chat/server/net/auth"
// 	mdwr "github.com/cubular-io/go-chat/server/net/middleware"
// 	"github.com/labstack/echo/v4"
// 	"github.com/labstack/echo/v4/middleware"
// )

// // Attaches a router with configurirable middleware and all the routes to the server
// func InitGroups(s *server.Server) {
// 	s.Router = &server.Router{
// 		// All Available Routes
// 		Routes: nil,

// 		Root: s.Echo.Group(""),

// 		Management: s.Echo.Group("/-", middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
// 			KeyLookup: "query:mgmt-secret",
// 			Validator: func(key string, c echo.Context) (bool, error) {
// 				return key == s.Config.Management.Secret, nil
// 			},
// 			Skipper: func(c echo.Context) bool {
// 				switch c.Path() {
// 				case "/-/ready":
// 					return true
// 				}
// 				return false
// 			},
// 		}), mdwr.NoCache()),

// 		V1User: s.Echo.Group("/v1/users", auth.WithConfig(auth.AuthConfig{
// 			S:              s,
// 			TokenValidator: auth.DefaultAuthTokenValidator,
// 			Mode:           auth.AuthModeRequired,
// 			Scopes:         []string{string(auth.AuthScopeAdmin)},
// 		})),

// 		V1Room: s.Echo.Group("/v1/rooms", auth.WithConfig(auth.AuthConfig{
// 			S:    s,
// 			Mode: auth.AuthModeNone,
// 		})),

// 		V1Msg: s.Echo.Group("/v1/msgs", auth.WithConfig(auth.AuthConfig{
// 			S:    s,
// 			Mode: auth.AuthModeNone,
// 		})),
// 	}
// }

// func AttachRoutes(s *server.Server) {
// 	// Attach all the routes
// 	s.Router.Routes = []*echo.Route{
// 		// == MANAGEMENT == //
// 		management.GetVersionRoute(s),
// 		management.GetDbVersionRoute(s),
// 		// == USER == //
// 		user.GetMeRoute(s),
// 		user.CreateUserRoute(s),
// 		user.UpdateUserRoute(s),
// 		user.DeleteUserRoute(s),
// 		// == ROOM == //
// 		room.CreateRoomRoute(s),
// 		// == MSG == //
// 		msg.CreateMsgRoute(s),
// 		msg.GetMsgRoute(s),
// 	}
// }
