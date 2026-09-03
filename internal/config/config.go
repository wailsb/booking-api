package config

import "os"

type Config struct {
	DBDSN string
	Port  string
}

func Load() *Config {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:secretpassword@localhost:5432/booking_db?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DBDSN: dsn,
		Port:  port,
	}
}
