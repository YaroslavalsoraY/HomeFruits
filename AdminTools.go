package main

import (
	"HomeFruits/internal/database"
	"HomeFruits/internal/jwt"
	"HomeFruits/logger"
	"context"
	"encoding/json"
	"net/http"
)

type Item struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Cost     int    `json:"cost"`
}

func (cfg *ApiConfig) HandlerInsertItem(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	adminID, err := jwt.ValidateJWT(token, cfg.SecretJWT)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	email, err := cfg.Queries.GetUserEmail(context.Background(), adminID)
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	if email != cfg.AdminEmail {
		http.Error(w, `{"error": "Do not have permissions"}`, http.StatusForbidden)
		logger.Warn(err)
		return
	}

	newItem := Item{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&newItem)
	if err != nil {
		http.Error(w, `{"error": "Problem with decoding json"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	args := database.InsertItemParams{
		Name:     newItem.Name,
		Quantity: int32(newItem.Quantity),
		Cost:     int32(newItem.Cost),
	}

	_, err = cfg.Queries.InsertItem(context.Background(), args)
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cfg *ApiConfig) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	adminID, err := jwt.ValidateJWT(token, cfg.SecretJWT)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	email, err := cfg.Queries.GetUserEmail(context.Background(), adminID)
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}
	if email != cfg.AdminEmail {
		http.Error(w, `{"error": "Do not have permissions"}`, http.StatusForbidden)
		logger.Warn(err)
		return
	}

	tokenToRevoke := r.PathValue("tokenID")

	cfg.Queries.RevokeToken(context.Background(), tokenToRevoke)

	w.WriteHeader(http.StatusNoContent)
}