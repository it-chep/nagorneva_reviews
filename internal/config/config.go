package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr, GRPCAddr, DatabaseURL, JWTSecret, S3Endpoint, S3Region, S3Bucket, S3AccessKey, S3SecretKey, S3PublicURL, AdminEmail, AdminPassword string
	PostgresHost, PostgresPort, PostgresUser, PostgresPassword, PostgresDB                                                                       string
	Debug                                                                                                                                        bool
	JWTTTL                                                                                                                                       time.Duration
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	ttl, err := time.ParseDuration(value("JWT_TTL", "24h"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_TTL: %w", err)
	}
	debug, err := strconv.ParseBool(value("DEBUG", "false"))
	if err != nil {
		return Config{}, fmt.Errorf("DEBUG: %w", err)
	}
	c := Config{
		HTTPAddr:         value("HTTP_ADDR", ":8080"),
		GRPCAddr:         value("GRPC_ADDR", ":7002"),
		PostgresHost:     value("POSTGRES_HOST", "localhost"),
		PostgresPort:     value("POSTGRES_PORT", "5432"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTTTL:           ttl,
		Debug:            debug,
		S3Endpoint:       value("S3_ENDPOINT", "https://storage.yandexcloud.net"),
		S3Region:         value("S3_REGION", "ru-central1"),
		S3Bucket:         os.Getenv("S3_BUCKET"),
		S3AccessKey:      os.Getenv("S3_ACCESS_KEY"),
		S3SecretKey:      os.Getenv("S3_SECRET_KEY"),
		S3PublicURL:      os.Getenv("S3_PUBLIC_URL"),
		AdminEmail:       os.Getenv("ADMIN_EMAIL"),
		AdminPassword:    os.Getenv("ADMIN_PASSWORD"),
	}
	if c.PostgresUser == "" || c.PostgresPassword == "" || c.PostgresDB == "" || len(c.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB and JWT_SECRET (at least 32 chars) are required")
	}
	c.DatabaseURL = postgresURL(c)
	return c, nil
}

func postgresURL(c Config) string {
	return (&url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.PostgresUser, c.PostgresPassword),
		Host:     net.JoinHostPort(c.PostgresHost, c.PostgresPort),
		Path:     c.PostgresDB,
		RawQuery: "sslmode=disable",
	}).String()
}

func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
