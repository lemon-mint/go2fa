package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/lemon-mint/envaddr"
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

	r := httprouter.New()

	r.POST("/api/v1/user/new", CreateUser)
	r.POST("/api/v1/totp/new", CreateTOTP)
	r.POST("/api/v1/totp/verify", VerifyTOTP)

	addr := envaddr.Get(":8080")
	log.Println("Listening on http://localhost" + addr)
	err = http.ListenAndServe(
		addr,
		r,
	)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln(err)
	}
}
