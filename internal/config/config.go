package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	HTTPPort     string
	LogBuffer    int
	ShutdownTime int
}

func Load() *Config {
	cfg := &Config{
		HTTPPort:     getEnv("HTTP_PORT", ":8080"),
		LogBuffer:    getEnvInt("LOG_BUFFER", 256),
		ShutdownTime: getEnvInt("SHUTDOWN_TIME", 10),
	}
	log.Printf("config loaded: %+v", cfg)
	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		var iv int
		_, err := fmt.Sscanf(v, "%d", &iv)
		if err == nil {
			return iv
		}
	}
	return def
}
