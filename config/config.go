package config

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all configuration values
type Config struct {
	Env          string
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPass     string
	KafkaBrokers []string
	APISecret    string
	NotifyEmail  string
}

// LoadConfig initializes Viper and returns the app config
func LoadConfig() (*Config, error) {
	// Detect environment (default: development)
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	// Map to the correct env file
	envFile := ".env"
	if env == "development" {
		envFile = ".env.dev"
	} else if env == "production" {
		envFile = ".env.prod"
	}

	viper.SetConfigFile(envFile)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("%s file not found. Falling back to environment variables.", envFile)
	} else {
		logrus.Infof("Loaded environment config from %s", envFile)
	}

	config := &Config{
		Env:          getEnv("ENV", env),
		Port:         getEnv("PORT", "8080"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "loan_service"),
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPass:     getEnv("SMTP_PASS", ""),
		KafkaBrokers: parseCSV("KAFKA_BROKERS"),
		APISecret:    getEnv("API_SECRET", "supersecretkey"),
		NotifyEmail:  getEnv("NOTIFY_EMAIL", "admin@example.com"),
	}

	return config, nil
}

func getEnv(key, fallback string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}
	return fallback
}

func parseCSV(key string) []string {
	raw := viper.GetString(key)
	if raw == "" {
		return []string{}
	}
	return strings.Split(raw, ",")
}