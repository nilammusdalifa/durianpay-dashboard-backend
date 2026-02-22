package config

import "os"

var (
	JwtSecret           = []byte(getEnv("JWT_SECRET", "dev-secret-replace-me"))
	JwtExpired          = getEnv("JWT_EXPIRED", "24h")
	HttpAddress         = getEnv("HTTP_ADDR", ":8080")
	OpenapiYamlLocation = getEnv("OPENAPIYAML_LOCATION", "../openapi.yaml")
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
