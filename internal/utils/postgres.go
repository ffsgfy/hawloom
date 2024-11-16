package utils

import (
	"fmt"
	"os"
)

func MakePostgresURI(user, password, host, port, db string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db)
}

func MakePostgresURIFromEnv() string {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	db := os.Getenv("POSTGRES_DB")
	if len(db) == 0 {
		db = user
	}

	return MakePostgresURI(user, password, host, port, db)
}
