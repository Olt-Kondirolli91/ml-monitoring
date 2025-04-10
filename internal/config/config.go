package config

import (
    "fmt"
    "os"
    "strconv"
)

type Config struct {
    DBHost     string
    DBPort     int
    DBUser     string
    DBPassword string
    DBName     string
    SSLMode    string
}

func LoadConfig() (*Config, error) {
    port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }

    return &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     port,
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "postgres"),
        SSLMode:    getEnv("DB_SSLMODE", "disable"),
    }, nil
}

func getEnv(key, defaultVal string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultVal
}
