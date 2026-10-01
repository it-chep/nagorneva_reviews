package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPAddr, DatabaseURL, JWTSecret, S3Endpoint, S3Region, S3Bucket, S3AccessKey, S3SecretKey, S3PublicURL, AdminEmail, AdminPassword string
	JWTTTL                                                                                                                             time.Duration
}

func Load() (Config, error) {
	ttl, err := time.ParseDuration(value("JWT_TTL", "24h"))
	if err != nil {
		return Config{}, fmt.Errorf("JWT_TTL: %w", err)
	}
	c := Config{HTTPAddr: value("HTTP_ADDR", ":8080"), DatabaseURL: os.Getenv("DATABASE_URL"), JWTSecret: os.Getenv("JWT_SECRET"), JWTTTL: ttl,
		S3Endpoint: os.Getenv("S3_ENDPOINT"), S3Region: value("S3_REGION", "us-east-1"), S3Bucket: os.Getenv("S3_BUCKET"), S3AccessKey: os.Getenv("S3_ACCESS_KEY"), S3SecretKey: os.Getenv("S3_SECRET_KEY"), S3PublicURL: os.Getenv("S3_PUBLIC_URL"), AdminEmail: os.Getenv("ADMIN_EMAIL"), AdminPassword: os.Getenv("ADMIN_PASSWORD")}
	if c.DatabaseURL == "" || len(c.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("DATABASE_URL and JWT_SECRET (at least 32 chars) are required")
	}
	return c, nil
}
func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
