package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	GrpcServerAddr      string
	RestPort            string
	AdminEmailAddr      string
	SenderEmailAddr     string
	SenderEmailPassword string
	LogLevel            string
	LogFile             string
	JWTSecret           string
	JWTExpiresIn        int
	RefreshTokenTTL     int
}

var (
	configInstance *Config
	once           sync.Once
)

func LoadConfig() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		jwtExpiresIn, err := strconv.Atoi(getEnv("JWT_EXPIRES_IN", "3600"))
		if err != nil {
			jwtExpiresIn = 3600
		}
		refreshTokenTTL, err := strconv.Atoi(getEnv("REFRESH_TOKEN_TTL", "604800"))
		if err != nil {
			refreshTokenTTL = 604800
		}

		configInstance = &Config{
			GrpcServerAddr:      getEnv("GRPC_SERVER_ADDRESS", "localhost:50051"),
			RestPort:            getEnv("REST_PORT", "8002"),
			AdminEmailAddr:      getEnv("ADMIN_EMAIL_ADDR", "winnerwinner2k3@gmail.com"),
			SenderEmailAddr:     getEnv("SENDER_EMAIL_ADDR", "thanhnt.works@gmail.com"),
			SenderEmailPassword: getEnv("SENDER_EMAIL_PASSWORD", "sdrgjqzoeosiklor"),
			LogLevel:            getEnv("LOG_LEVEL", "info"),
			LogFile:             getEnv("LOG_FILE", "../logs/mail.log"),
			JWTSecret:           getEnv("JWT_SECRET", "supersecretkey"),
			JWTExpiresIn:        jwtExpiresIn,
			RefreshTokenTTL:     refreshTokenTTL,
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
