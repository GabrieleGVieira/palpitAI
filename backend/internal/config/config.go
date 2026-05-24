package config

import "os"

type Config struct {
	DatabaseURL                 string
	Env                         string
	FootballDataAPIBaseURL      string
	FootballDataCompetitionCode string
	FootballDataSeason          string
	FootballDataToken           string
	Port                        string
	RedisURL                    string
	SupabaseKey                 string
	SupabaseURL                 string
}

func Load() Config {
	return Config{
		DatabaseURL:                 getEnv("DATABASE_URL", "postgresql://postgres:jpEWsj0V4iT2VRaz@db.kqczvvcnlhpukrrioctp.supabase.co:5432/postgres"),
		Env:                         getEnv("APP_ENV", "development"),
		FootballDataAPIBaseURL:      getEnv("FOOTBALL_DATA_API_BASE_URL", "https://api.football-data.org/v4"),
		FootballDataCompetitionCode: getEnv("FOOTBALL_DATA_COMPETITION_CODE", "BSA"),
		FootballDataSeason:          getEnv("FOOTBALL_DATA_SEASON", ""),
		FootballDataToken:           getEnv("FOOTBALL_DATA_TOKEN", "2401e3dee22043c1acefc19ba5620165"),
		Port:                        getEnv("PORT", "3000"),
		RedisURL:                    getEnv("REDIS_URL", "redis://default:gQAAAAAAAfpqAAIgcDJmYTlmNjE1YmYxY2E0YTAwODYwYzZjZjFhMDVjNTJjNw@amusing-halibut-129642.upstash.io:6379"),
		SupabaseKey:                 getEnv("SUPABASE_KEY", ""),
		SupabaseURL:                 getEnv("SUPABASE_URL", ""),
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
