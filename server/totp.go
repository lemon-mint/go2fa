package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base32"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/lemon-mint/go2fa/methods/totp"
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

func VerifyTOTP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type VerifyRequest struct {
		UserID uint64 `json:"user_id"`
		Token  string `json:"token"`
	}

	var verifyRequest VerifyRequest
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&verifyRequest)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
	} else if r.Method == "GET" {
		userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 64)
		if err != nil {
			http.Error(w, "Bad Request", 400)
			return
		}
		verifyRequest.UserID = userID
		verifyRequest.Token = r.URL.Query().Get("token")
	}
	var keyID uint64
	authMethodsRows, err := DB.Query(
		context.Background(),
		`SELECT data_id FROM auth_methods WHERE user_id = $1 AND type = $2`,
		verifyRequest.UserID,
		Auth_Method_TOTP,
	)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	defer authMethodsRows.Close()
	for authMethodsRows.Next() {
		err = authMethodsRows.Scan(&keyID)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		var secret string
		err = DB.QueryRow(
			context.Background(),
			`SELECT secret FROM totp_tokens WHERE id = $1`,
			keyID,
		).Scan(&secret)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		secretBytes, err := base32.StdEncoding.DecodeString(secret)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}

		otp := totp.TOTP{
			Secret: secretBytes,
			Digits: 6,
			Period: 30,
		}

		serverToken := otp.Generate(time.Now())

		// Securely compare the two strings.
		if subtle.ConstantTimeCompare([]byte(serverToken), []byte(verifyRequest.Token)) == 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(`{"success": true}`))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{"success": false}`))
}
