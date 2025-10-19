package config

import (
	"os"
	"time"
)

type Config struct {
	Server    ServerConfig
	Scraper   ScraperConfig
	Scheduler SchedulerConfig
	Storage   StorageConfig
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type ScraperConfig struct {
	BaseURL        string
	RequestTimeout time.Duration
	RateLimit      time.Duration
	UserAgent      string
	AcceptLanguage string
}

type SchedulerConfig struct {
	CronExpression string
	Timezone       string
}

type StorageConfig struct {
	FilePath string
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     getDuration("READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getDuration("WRITE_TIMEOUT", 15*time.Second),
			ShutdownTimeout: getDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Scraper: ScraperConfig{
			BaseURL:        getEnv("IMDB_BASE_URL", "https://www.imdb.com"),
			RequestTimeout: getDuration("REQUEST_TIMEOUT", 30*time.Second),
			RateLimit:      getDuration("RATE_LIMIT", 2*time.Second),
			UserAgent:      getEnv("USER_AGENT", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
			AcceptLanguage: getEnv("ACCEPT_LANGUAGE", "en-US,en;q=0.9"),
		},
		Scheduler: SchedulerConfig{
			CronExpression: getEnv("CRON_SCHEDULE", "0 9 * * *"),
			Timezone:       getEnv("TIMEZONE", "UTC"),
		},
		Storage: StorageConfig{
			FilePath: getEnv("STORAGE_FILE", "data/tracked_series.json"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
