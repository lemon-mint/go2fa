package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user_id uint64
	err := DB.QueryRow(
		context.Background(),
		`INSERT INTO users (created_at, updated_at) VALUES (now(), now()) RETURNING user_id`,
	).Scan(&user_id)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}

	type UserInfo struct {
		Success bool   `json:"success"`
		UserID  uint64 `json:"user_id"`
	}

	userInfo := UserInfo{
		Success: true,
		UserID:  user_id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(userInfo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
