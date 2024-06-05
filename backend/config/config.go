package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Conf struct {
	Server   ConfServer
	Database ConfDatabase
}

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
	APIKey       string        `env:"SERVER_API_KEY,required"`
	Secret       string        `env:"SERVER_SECRET,required"`
}

type ConfDatabase struct {
	Host     string `env:"POSTGRES_HOST,required"`
	Port     int    `env:"POSTGRES_PORT,required"`
	Username string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	Name     string `env:"POSTGRES_DB,required"`
	Debug    bool   `env:"POSTGRES_DEBUG,required"`
}

func New() *Conf {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env: %s", err)
	}

	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}

func NewDatabase() *ConfDatabase {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load env: %s", err)
	}

	var c ConfDatabase
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
