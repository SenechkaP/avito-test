package configs

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerConfig ServerConfig
	DbConfig     DbConfig
}

type ServerConfig struct {
	Port string
}

type DbConfig struct {
	User              string
	Password          string
	Name              string
	Host              string
	Port              string
	AttemptsToConnect int
	BaseDelay         time.Duration
}

func (db *DbConfig) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db.User, db.Password, db.Host, db.Port, db.Name)
}

func LoadConfig(envPath string) *Config {
	if envPath != "" {
		if err := godotenv.Load(envPath); err != nil {
			log.Fatalf("Could not load env file %s: %v", envPath, err)
		}
	}

	attempts, _ := strconv.Atoi(getEnv("ATTEMPTS_TO_CONNECT", "5"))
	baseDelay, _ := strconv.Atoi(getEnv("BASE_DELAY", "2"))

	return &Config{
		ServerConfig: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		DbConfig: DbConfig{
			User:              getEnv("POSTGRES_USER", "user"),
			Password:          getEnv("POSTGRES_PASSWORD", "pass"),
			Name:              getEnv("POSTGRES_DB_NAME", "test_db"),
			Host:              getEnv("POSTGRES_HOST", "localhost"),
			Port:              getEnv("POSTGRES_PORT", "5432"),
			AttemptsToConnect: attempts,
			BaseDelay:         time.Second * time.Duration(baseDelay),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}
