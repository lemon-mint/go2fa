package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lemon-mint/godotenv"
)

var DB *pgxpool.Pool

func main() {
	var err error
	godotenv.Load()
	DB, err = pgxpool.Connect(
		context.Background(),
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer DB.Close()

	err = InitDB()
	if err != nil {
		log.Fatalln(err)
	}

}
