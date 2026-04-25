package shared

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JwtSecret []byte
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Print("no .env file, reading from environment")
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	return &Config{JwtSecret: jwtSecret}
}
