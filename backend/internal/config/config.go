package config

import "os"

type Config struct {
	DatabaseURL string
	Env         string
	Port        string
	SupabaseKey string
	SupabaseURL string
}

func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Env:         getEnv("APP_ENV", "development"),
		Port:        getEnv("PORT", "3000"),
		SupabaseKey: getEnv("SUPABASE_KEY", ""),
		SupabaseURL: getEnv("SUPABASE_URL", ""),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
