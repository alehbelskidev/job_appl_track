package shared

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JwtSecret   []byte
	DatabaseUrl string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file, reading from environment")
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	databaseUrl := os.Getenv("DATABASE_URL")

	return &Config{JwtSecret: jwtSecret, DatabaseUrl: databaseUrl}
}
