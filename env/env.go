package env

import (
	"fmt"
	"os"
)

var (
	Port      = getEnv("PORT", "8000")
	CoreDBDSN = getEnvOrPanic("CORE_DB_DSN")

	JWTSigningKeyARN = getEnvOrPanic("JWT_SIGNING_KEY_ARN")
	JWTSigningKeyKID = getEnvOrPanic("JWT_SIGNING_KEY_KID")
)

func getEnv(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return fallback
}

func getEnvOrPanic(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("❌ missing required environment variable: '%v'\n", key))
	}
	return value
}
