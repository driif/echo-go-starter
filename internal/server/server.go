package server

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strings"

	"github.com/driif/echo-go-starter/internal/server/config"
	mdwr "github.com/driif/echo-go-starter/internal/server/net/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

// Server is the main struct for the server
type Server struct {
	// Code here
	Config config.Server
	Echo   *echo.Echo
	Router *Router
	DB     *sql.DB
	//Mailer *mailer.Mailer
	//Push *push.Service
}

// Router is a struct that holds all the routes for the server
type Router struct {
	Routes     []*echo.Route
	Root       *echo.Group
	Management *echo.Group
	V1User     *echo.Group
	V1Room     *echo.Group
	V1Msg      *echo.Group
}

// New creates a new server
func New(config *config.Server) *Server {
	// Code here
	s := &Server{
		Config: *config,
		DB:     nil,
		Echo:   nil,
		Router: nil,
		//Router: nil,
	}
	return s
}

// Checks if the server is ready to serve requests
func (s *Server) Ready() bool {
	return s.DB != nil &&
		s.Echo != nil &&
		s.Router != nil
}

// InitPush initializes the push service
func (s *Server) InitPush() error {
	// Code here
	return nil
}

// Initialize a new Echo server with Middleware Configs
func (s *Server) Initialize() error {
	s.Echo = echo.New()

	s.Echo.Debug = s.Config.Echo.Debug
	s.Echo.HideBanner = true
	s.Echo.Logger.SetOutput(&echoLogger{level: s.Config.Logger.RequestLevel, log: log.With().Str("component", "echo").Logger()})

	// add handler before each route
	s.Echo.HTTPErrorHandler = HTTPErrorHandlerWithConfig(HTTPErrorHandlerConfig{
		HideInternalServerErrorDetails: s.Config.Echo.HideInternalServerErrorDetails,
	})

	// ---
	// General middleware
	if s.Config.Echo.EnableTrailingSlashMiddleware {
		s.Echo.Pre(middleware.RemoveTrailingSlash())
	} else {
		log.Warn().Msg("Disabling trailing slash middleware due to environment config")
	}

	if s.Config.Echo.EnableTrailingSlashMiddleware {
		s.Echo.Use(middleware.Recover())
	} else {
		log.Warn().Msg("Disabling recover middleware due to environment config")
	}

	if s.Config.Echo.EnableSecureMiddleware {
		s.Echo.Use(middleware.SecureWithConfig(middleware.SecureConfig{
			Skipper:               middleware.DefaultSecureConfig.Skipper,
			XSSProtection:         s.Config.Echo.SecureMiddleware.XSSProtection,
			ContentTypeNosniff:    s.Config.Echo.SecureMiddleware.ContentTypeNosniff,
			XFrameOptions:         s.Config.Echo.SecureMiddleware.XFrameOptions,
			HSTSMaxAge:            s.Config.Echo.SecureMiddleware.HSTSMaxAge,
			HSTSExcludeSubdomains: s.Config.Echo.SecureMiddleware.HSTSExcludeSubdomains,
			ContentSecurityPolicy: s.Config.Echo.SecureMiddleware.ContentSecurityPolicy,
			CSPReportOnly:         s.Config.Echo.SecureMiddleware.CSPReportOnly,
			HSTSPreloadEnabled:    s.Config.Echo.SecureMiddleware.HSTSPreloadEnabled,
			ReferrerPolicy:        s.Config.Echo.SecureMiddleware.ReferrerPolicy,
		}))
	} else {
		log.Warn().Msg("Disabling secure middleware due to environment config")
	}

	if s.Config.Echo.EnableCORSMiddleware {
		s.Echo.Use(middleware.RequestID())
	} else {
		log.Warn().Msg("Disabling request ID middleware due to environment config")
	}

	if s.Config.Echo.EnableLoggerMiddleware {
		s.Echo.Use(mdwr.LoggerWithConfig(mdwr.LoggerConfig{
			Level:             s.Config.Logger.RequestLevel,
			LogRequestBody:    s.Config.Logger.LogRequestBody,
			LogRequestHeader:  s.Config.Logger.LogRequestHeader,
			LogRequestQuery:   s.Config.Logger.LogRequestQuery,
			LogResponseBody:   s.Config.Logger.LogResponseBody,
			LogResponseHeader: s.Config.Logger.LogResponseHeader,
			RequestBodyLogSkipper: func(req *http.Request) bool {
				// Skip all body logging for auth endpoints as these might contain sensitive data
				if strings.HasPrefix(req.URL.Path, "/v1/auth") {
					return true
				}

				return mdwr.DefaultRequestBodyLogSkipper(req)
			},
			ResponseBodyLogSkipper: func(req *http.Request, res *echo.Response) bool {
				// We skip all body logging for auth endpoints as these might contain sensitive data
				if strings.HasPrefix(req.URL.Path, "/v1/auth") {
					return true
				}

				return mdwr.DefaultResponseBodyLogSkipper(req, res)
			},
			Skipper: func(c echo.Context) bool {
				// We skip loggging of readiness and liveness endpoints
				switch c.Path() {
				case "/-/ready", "/-/healthy":
					return true
				}
				return false
			},
		}))
	} else {
		log.Warn().Msg("Disabling logger middleware due to environment config")
	}

	if s.Config.Echo.EnableCORSMiddleware {
		s.Echo.Use(middleware.CORS())
	} else {
		log.Warn().Msg("Disabling CORS middleware due to environment config")
	}

	if s.Config.Echo.EnableCacheControlMiddleware {
		s.Echo.Use(mdwr.CacheControl())
	} else {
		log.Warn().Msg("Disabling cache control middleware due to environment config")
	}

	if s.Config.Pprof.Enable {
		pprofAuthMiddleware := mdwr.Noop()

		if s.Config.Pprof.EnableManagementKeyAuth {
			pprofAuthMiddleware = middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
				KeyLookup: "query:mgmt-secret",
				Validator: func(key string, c echo.Context) (bool, error) {
					return key == s.Config.Management.Secret, nil
				},
			})
		}

		s.Echo.GET("/debug/pprof", echo.WrapHandler(http.HandlerFunc(pprof.Index)), pprofAuthMiddleware)
		s.Echo.Any("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux), pprofAuthMiddleware)

		log.Warn().Bool("EnableManagementKeyAuth", s.Config.Pprof.EnableManagementKeyAuth).Msg("Pprof http handlers are avaible at /debug/pprof")

		if s.Config.Pprof.RuntimeMutexProfileFraction != 0 {
			runtime.SetMutexProfileFraction(s.Config.Pprof.RuntimeMutexProfileFraction)
			log.Warn().Int("RuntimeMutexProfileFraction", s.Config.Pprof.RuntimeMutexProfileFraction).Msg("Pprof runtime.SetMutexProfileFraction")
		}
	}

	// Code here
	return nil
}

// Starts the server
func (s *Server) Start() error {
	// Code here
	return s.Echo.Start(s.Config.Echo.ListenAddress)
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {

	// Code here
	log.Warn().Msg("Shutting down server")

	return s.Echo.Shutdown(ctx)
}
