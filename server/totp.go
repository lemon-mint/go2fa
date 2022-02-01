package main

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
)

const (
	Auth_Method_TOTP = 1
)

type User struct {
	UserID uint64 `json:"user_id"`
}

func CreateTOTP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	secret := make([]byte, 20)
	_, err = rand.Read(secret)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	secretStr := base32.StdEncoding.EncodeToString(secret)
	var keyID uint64
	err = DB.BeginFunc(
		context.Background(),
		func(t pgx.Tx) error {
			err = t.QueryRow(
				context.Background(),
				`INSERT INTO totp_tokens (secret) VALUES ($1) RETURNING id`,
				secretStr,
			).Scan(&keyID)
			if err != nil {
				return err
			}
			_, err = t.Exec(
				context.Background(),
				`INSERT INTO auth_methods (user_id, type, data_id) VALUES ($1, $2, $3)`,
				user.UserID,
				Auth_Method_TOTP,
				keyID,
			)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	type TOTPInfo struct {
		Success bool   `json:"success"`
		Secret  string `json:"secret"`
	}

	totpInfo := TOTPInfo{
		Success: true,
		Secret:  secretStr,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(totpInfo)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
}
