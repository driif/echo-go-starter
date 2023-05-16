package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/driif/echo-go-starter/internal/server/config/env"
	"github.com/driif/echo-go-starter/pkg/logs"
	"github.com/driif/echo-go-starter/pkg/tests"
	"github.com/rs/zerolog"
)

// EchoServer represents a subset of echo's config relevant to the app server.
type EchoServer struct {
	Debug                          bool
	ListenAddress                  string
	HideInternalServerErrorDetails bool
	BaseURL                        string
	EnableCORSMiddleware           bool
	EnableLoggerMiddleware         bool
	EnableRecoverMiddleware        bool
	EnableRequestIDMiddleware      bool
	EnableTrailingSlashMiddleware  bool
	EnableSecureMiddleware         bool
	EnableCacheControlMiddleware   bool
	SecureMiddleware               EchoServerSecureMiddleware
}

// PprofServer represents a subset of pprof's config relevant to the app server.
type PprofServer struct {
	Enable                      bool
	EnableManagementKeyAuth     bool
	RuntimeBlockProfileRate     int
	RuntimeMutexProfileFraction int
}

// EchoServerSecureMiddleware represents a subset of echo's secure middleware config relevant to the app server.
// https://github.com/labstack/echo/blob/master/middleware/secure.go
type EchoServerSecureMiddleware struct {
	XSSProtection         string
	ContentTypeNosniff    string
	XFrameOptions         string
	HSTSMaxAge            int
	HSTSExcludeSubdomains bool
	ContentSecurityPolicy string
	CSPReportOnly         bool
	HSTSPreloadEnabled    bool
	ReferrerPolicy        string
}

// PathsServer represents a subset of paths config relevant to the app server.
type PathsServer struct {
	APIBaseDirAbs string
	MntBaseDirAbs string
}

// ManagementServer represents a subset of management config relevant to the app server.
type ManagementServer struct {
	Secret                  string `json:"-"` // sensitive
	CryptoKey               string `json:"-"` // sensitive
	ReadinessTimeout        time.Duration
	LivenessTimeout         time.Duration
	ProbeWriteablePathsAbs  []string
	ProbeWriteableTouchfile string
}

// LoggerServer represents a subset of logger config relevant to the app server.
type LoggerServer struct {
	Level              zerolog.Level
	RequestLevel       zerolog.Level
	LogRequestBody     bool
	LogRequestHeader   bool
	LogRequestQuery    bool
	LogResponseBody    bool
	LogResponseHeader  bool
	PrettyPrintConsole bool
}

// Server represents the config of the Server relevant to the app server, containing all the other config structs.
type Server struct {
	Database   Database
	Echo       EchoServer
	Pprof      PprofServer
	Paths      PathsServer
	Management ManagementServer
	//Mailer     Mailer
	//SMTP       transport.SMTPMailTransportConfig
	//Frontend   FrontendServer
	Logger LoggerServer
	//Push       PushService
	//FCMConfig  provider.FCMConfig
}

// DefaultServiceConfigFromEnv returns the server config as parsed from environment variables
// and their respective defaults defined below.
// We don't expect that ENV_VARs change while we are running our application or our tests
// (and it would be a bad thing to do anyways with parallel testing).
// Do NOT use os.Setenv / os.Unsetenv in tests envizing DefaultServiceConfigFromEnv()!
func DefaultServiceConfigFromEnv() Server {

	// An `.env.local` file in your project root can override the currently set ENV variables.
	//
	// We never automatically apply `.env.local` when running "go test" as these ENV variables
	// may be sensitive (e.g. secrets to external APIs) and applying them modifies the process
	// global "os.Env" state (it should be applied via t.SetEnv instead).
	//
	// If you need dotenv ENV variables available in a test, do that explicitly within that
	// test before executing DefaultServiceConfigFromEnv (or test.WithTestServer).
	// See /internal/test/helper_dot_env.go: test.DotEnvLoadLocalOrSkipTest(t)
	if !tests.RunningInTest() {
		env.DotEnvTryLoad(filepath.Join(env.GetProjectRootDir(), ".env.local"), os.Setenv)
	}

	return Server{
		Database: Database{
			Host:     env.GetEnv("PGHOST", "postgres"),
			Port:     env.GetEnvAsInt("PGPORT", 5432),
			Database: env.GetEnv("PGDATABASE", "development"),
			Username: env.GetEnv("PGUSER", "dbuser"),
			Password: env.GetEnv("PGPASSWORD", "dbpass"),
			AdditionalParams: map[string]string{
				"sslmode": env.GetEnv("PGSSLMODE", "disable"),
			},
			MaxOpenConns:    env.GetEnvAsInt("DB_MAX_OPEN_CONNS", runtime.NumCPU()*2),
			MaxIdleConns:    env.GetEnvAsInt("DB_MAX_IDLE_CONNS", 1),
			ConnMaxLifetime: time.Second * time.Duration(env.GetEnvAsInt("DB_CONN_MAX_LIFETIME_SEC", 60)),
		},
		Echo: EchoServer{
			Debug:                          env.GetEnvAsBool("SERVER_ECHO_DEBUG", false),
			ListenAddress:                  env.GetEnv("SERVER_ECHO_LISTEN_ADDRESS", ":8080"),
			HideInternalServerErrorDetails: env.GetEnvAsBool("SERVER_ECHO_HIDE_INTERNAL_SERVER_ERROR_DETAILS", true),
			BaseURL:                        env.GetEnv("SERVER_ECHO_BASE_URL", "http://localhost:8080"),
			EnableCORSMiddleware:           env.GetEnvAsBool("SERVER_ECHO_ENABLE_CORS_MIDDLEWARE", true),
			EnableLoggerMiddleware:         env.GetEnvAsBool("SERVER_ECHO_ENABLE_LOGGER_MIDDLEWARE", true),
			EnableRecoverMiddleware:        env.GetEnvAsBool("SERVER_ECHO_ENABLE_RECOVER_MIDDLEWARE", true),
			EnableRequestIDMiddleware:      env.GetEnvAsBool("SERVER_ECHO_ENABLE_REQUEST_ID_MIDDLEWARE", true),
			EnableTrailingSlashMiddleware:  env.GetEnvAsBool("SERVER_ECHO_ENABLE_TRAILING_SLASH_MIDDLEWARE", true),
			EnableSecureMiddleware:         env.GetEnvAsBool("SERVER_ECHO_ENABLE_SECURE_MIDDLEWARE", true),
			EnableCacheControlMiddleware:   env.GetEnvAsBool("SERVER_ECHO_ENABLE_CACHE_CONTROL_MIDDLEWARE", true),
			// see https://echo.labstack.com/middleware/secure
			// see https://github.com/labstack/echo/blob/master/middleware/secure.go
			SecureMiddleware: EchoServerSecureMiddleware{
				XSSProtection:         env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_XSS_PROTECTION", "1; mode=block"),
				ContentTypeNosniff:    env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_CONTENT_TYPE_NOSNIFF", "nosniff"),
				XFrameOptions:         env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_X_FRAME_OPTIONS", "SAMEORIGIN"),
				HSTSMaxAge:            env.GetEnvAsInt("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_MAX_AGE", 0),
				HSTSExcludeSubdomains: env.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_EXCLUDE_SUBDOMAINS", false),
				ContentSecurityPolicy: env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_CONTENT_SECURITY_POLICY", ""),
				CSPReportOnly:         env.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_CSP_REPORT_ONLY", false),
				HSTSPreloadEnabled:    env.GetEnvAsBool("SERVER_ECHO_SECURE_MIDDLEWARE_HSTS_PRELOAD_ENABLED", false),
				ReferrerPolicy:        env.GetEnv("SERVER_ECHO_SECURE_MIDDLEWARE_REFERRER_POLICY", ""),
			},
		},
		Pprof: PprofServer{
			// https://golang.org/pkg/net/http/pprof/
			Enable:                      env.GetEnvAsBool("SERVER_PPROF_ENABLE", false),
			EnableManagementKeyAuth:     env.GetEnvAsBool("SERVER_PPROF_ENABLE_MANAGEMENT_KEY_AUTH", true),
			RuntimeBlockProfileRate:     env.GetEnvAsInt("SERVER_PPROF_RUNTIME_BLOCK_PROFILE_RATE", 0),
			RuntimeMutexProfileFraction: env.GetEnvAsInt("SERVER_PPROF_RUNTIME_MUTEX_PROFILE_FRACTION", 0),
		},
		Paths: PathsServer{
			// Please ALWAYS work with ABSOLUTE (ABS) paths from ENV_VARS (however you may resolve a project-relative to absolute for the default value)
			APIBaseDirAbs: env.GetEnv("SERVER_PATHS_API_BASE_DIR_ABS", filepath.Join(env.GetProjectRootDir(), "/api")),        // /app/api (swagger.yml)
			MntBaseDirAbs: env.GetEnv("SERVER_PATHS_MNT_BASE_DIR_ABS", filepath.Join(env.GetProjectRootDir(), "/assets/mnt")), // /app/assets/mnt (user-generated content)
		},
		Management: ManagementServer{
			Secret:           env.GetMgmtSecret("SERVER_MANAGEMENT_SECRET"),
			CryptoKey:        env.GetEnv("CRYPTO_KEY", "12345678901234567-Ia123456789012"),
			ReadinessTimeout: time.Second * time.Duration(env.GetEnvAsInt("SERVER_MANAGEMENT_READINESS_TIMEOUT_SEC", 4)),
			LivenessTimeout:  time.Second * time.Duration(env.GetEnvAsInt("SERVER_MANAGEMENT_LIVENESS_TIMEOUT_SEC", 9)),
			ProbeWriteablePathsAbs: env.GetEnvAsStringArr("SERVER_MANAGEMENT_PROBE_WRITEABLE_PATHS_ABS", []string{
				filepath.Join(env.GetProjectRootDir(), "/assets/mnt")}, ","),
			ProbeWriteableTouchfile: env.GetEnv("SERVER_MANAGEMENT_PROBE_WRITEABLE_TOUCHFILE", ".healthy"),
		},
		Logger: LoggerServer{
			Level:              logs.LogLevelFromString(env.GetEnv("SERVER_LOGGER_LEVEL", zerolog.DebugLevel.String())),
			RequestLevel:       logs.LogLevelFromString(env.GetEnv("SERVER_LOGGER_REQUEST_LEVEL", zerolog.DebugLevel.String())),
			LogRequestBody:     env.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_BODY", false),
			LogRequestHeader:   env.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_HEADER", false),
			LogRequestQuery:    env.GetEnvAsBool("SERVER_LOGGER_LOG_REQUEST_QUERY", false),
			LogResponseBody:    env.GetEnvAsBool("SERVER_LOGGER_LOG_RESPONSE_BODY", false),
			LogResponseHeader:  env.GetEnvAsBool("SERVER_LOGGER_LOG_RESPONSE_HEADER", false),
			PrettyPrintConsole: env.GetEnvAsBool("SERVER_LOGGER_PRETTY_PRINT_CONSOLE", false),
		},
	}

}
