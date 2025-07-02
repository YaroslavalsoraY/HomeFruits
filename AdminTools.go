package main

import (
	"HomeFruits/internal/database"
	"HomeFruits/internal/jwt"
	"HomeFruits/logger"
	"context"
	"encoding/json"
	"fmt"
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
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	adminID, err := jwt.ValidateJWT(token, cfg.SecretJWT)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	email, err := cfg.Queries.GetUserEmail(context.Background(), adminID)
	if err != nil {
		fmt.Println(adminID)
		logger.Warn(err)
		return
	}

	if email != cfg.AdminEmail {
		w.WriteHeader(http.StatusUnavailableForLegalReasons)
		logger.Warn(err)
		return
	}

	newItem := Item{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&newItem)

	args := database.InsertItemParams{
		Name:     newItem.Name,
		Quantity: int32(newItem.Quantity),
		Cost:     int32(newItem.Cost),
	}

	_, err = cfg.Queries.InsertItem(context.Background(), args)
	if err != nil {
		logger.Warn(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
