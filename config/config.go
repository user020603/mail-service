package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	GrpcServerAddr      string
	AdminEmailAddr      string
	SenderEmailAddr     string
	SenderEmailPassword string
	LogLevel            string
	LogFile             string
}

var (
	configInstance *Config
	once           sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		configInstance = &Config{
			GrpcServerAddr:      getEnv("GRPC_SERVER_ADDRESS", "localhost:50051"),
			AdminEmailAddr:      getEnv("ADMIN_EMAIL_ADDRESS", ""),
			SenderEmailAddr:     getEnv("SENDER_EMAIL_ADDRESS", ""),
			SenderEmailPassword: getEnv("SENDER_EMAIL_PASSWORD", ""),
			LogLevel:            getEnv("LOG_LEVEL", "info"),
			LogFile:             getEnv("LOG_FILE", "../logs/mail.log"),
		}
	})

	return configInstance
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
