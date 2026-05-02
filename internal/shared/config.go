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
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file, reading from environment")
	}

	allowedOriginsStr := []byte(os.Getenv("ALLOWED_ORIGINS"))
	allowedOrigins := strings.Split(string(allowedOriginsStr), ",")

	pgUser := []byte(os.Getenv("POSTGRES_USER"))
	pgPassword := []byte(os.Getenv("POSTGRES_PASSWORD"))
	pgDB := []byte(os.Getenv("POSTGRES_DB"))
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	databaseUrl := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", pgUser, pgPassword, pgDB)
	databaseUrlMigrate := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", pgUser, pgPassword, pgDB)

	return &Config{
		JwtSecret:          jwtSecret,
		DatabaseUrl:        databaseUrl,
		DatabaseUrlMigrate: databaseUrlMigrate,
		AllowedOrigins:     allowedOrigins,
	}
}
