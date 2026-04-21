package config

import "os"

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    ServerPort string
}

func LoadConfig() *Config {
    return &Config{
        DBHost:     getEnv("DB_HOST", getEnv("PGHOST", "localhost")),
        DBPort:     getEnv("DB_PORT", getEnv("PGPORT", "5432")),
        DBUser:     getEnv("DB_USER", getEnv("PGUSER", "postgres")),
        DBPassword: getEnv("DB_PASSWORD", getEnv("PGPASSWORD", "")),
        DBName:     getEnv("DB_NAME", getEnv("PGDATABASE", "happypc")),
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}