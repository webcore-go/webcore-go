package config

func (c *Config) SetDefaults() map[string]any {
	return map[string]any{
		// App
		"app.name":                   "webcore-go",
		"app.version":                "1.0.0",
		"app.environment":            "development",
		"app.features.recovery":      false,
		"app.features.tracing":       false,
		"app.features.metrics":       false,
		"app.features.profiling":     false,
		"app.logging.level":          "info",
		"app.logging.format":         "json",
		"app.logging.output":         "stdout",
		"app.security_headers":       false,
		"app.cors.allow_origins":     []string{"*"},
		"app.cors.allow_methods":     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		"app.cors.allow_headers":     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		"app.cors.expose_headers":    []string{"Content-Length"},
		"app.cors.allow_credentials": true,
		"app.cors.max_age":           "24h", // 24 hours
		"app.rate_limit.enabled":     false,
		"app.rate_limit.max":         1000,
		"app.module.base_path":       "./libs",
		"app.module.disabled":        []string{},

		// Server
		"server.host":          "0.0.0.0",
		"server.port":          7272,
		"server.path":          "/api",
		"server.read_timeout":  "30s",
		"server.write_timeout": "30s",

		// Auth
		"auth.type":           "jwt",
		"auth.secret_key":     "",
		"auth.expires_in":     "24h", // 24 hours
		"auth.api_key_header": "X-API-Key",
		"auth.api_key_prefix": "",

		// Database
		"database.driver":            "postgres",
		"database.host":              "",
		"database.port":              5432,
		"database.ssl_mode":          "disable",
		"database.max_idle_conns":    10,
		"database.max_open_conns":    100,
		"database.conn_max_lifetime": "300s", // in seconds

		// Redis
		"redis.host": "",
		"redis.port": 6379,
		"redis.db":   0,

		// PubSub
		"pubsub.project_id":               "",
		"pubsub.topic":                    "",
		"pubsub.subscription":             "",
		"pubsub.emulator_host":            "",
		"pubsub.credentials":              "",
		"pubsub.max_messages_per_pull":    50,
		"pubsub.sleep_time_between_pulls": "10s", // in seconds
	}
}
