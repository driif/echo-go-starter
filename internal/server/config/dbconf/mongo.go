package dbconf

import "github.com/driif/echo-go-starter/internal/server/config/env"

// MongoDb - config for MongoDB
type MongoDB struct {
	URI        string
	Host       string
	Port       int
	Database   string
	Username   string
	Password   string
	MaxResults int
	// This is capped by the Session's send queue limit (128).
	MaxMessageResults int
	AdpVersion        int
	AdapterName       string
}

func ConfigMongo() *MongoDB {
	return &MongoDB{
		URI:               env.GetEnv("MG_URI", "mongodb://admin:pass@localhost:27017"),
		Host:              env.GetEnv("MG_HOST", "localhost"),
		Port:              env.GetEnvAsInt("MG_PORT", 27017),
		Database:          env.GetEnv("MG_DATABASE", "chat"),
		Username:          env.GetEnv("MG_USERNAME", "chat"),
		Password:          env.GetEnv("MG_PASSWORD", "chat"),
		MaxResults:        env.GetEnvAsInt("MG_MAX_RESULTS", 1024),
		MaxMessageResults: env.GetEnvAsInt("MG_MAX_MESSAGE_RESULTS", 100),
		AdpVersion:        env.GetEnvAsInt("MG_ADAPTER_VERSION", 112),
		AdapterName:       env.GetEnv("MG_ADAPTER_NAME", "mongodb"),
	}
}
