package shared

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	JwtSecret          []byte
	DatabaseUrl        string
	DatabaseUrlMigrate string
	AllowedOrigins     []string
	OpenAISecret       []byte
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file, reading from environment")
	}

	allowedOriginsStr := []byte(os.Getenv("ALLOWED_ORIGINS"))
	allowedOrigins := strings.Split(string(allowedOriginsStr), ",")

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	pgUser := []byte(os.Getenv("POSTGRES_USER"))
	pgPassword := []byte(os.Getenv("POSTGRES_PASSWORD"))
	pgDB := []byte(os.Getenv("POSTGRES_DB"))
	pgHost := []byte(os.Getenv("POSTGRES_HOST"))
	pgHostMigrate := []byte(os.Getenv("POSTGRES_HOST_MIGRATE"))
	openAISecret := []byte(os.Getenv("OPENAI_SECRET"))

	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		pgUser, pgPassword, pgHost, pgDB,
	)
	databaseUrlMigrate := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/%s?sslmode=disable",
		pgUser, pgPassword, pgHostMigrate, pgDB,
	)

	return &Config{
		JwtSecret:          jwtSecret,
		DatabaseUrl:        databaseUrl,
		DatabaseUrlMigrate: databaseUrlMigrate,
		AllowedOrigins:     allowedOrigins,
		OpenAISecret:       openAISecret,
	}
}
