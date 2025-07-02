package main

import (
	"HomeFruits/internal/database"
	"HomeFruits/internal/jwt"
	"HomeFruits/logger"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type GetItemParams struct {
	ItemID   uuid.UUID
	UserID   uuid.UUID
	Name     string
	Quantity int `json:"quantity"`
	Cost     int `json:"cost"`
}

func (cfg *ApiConfig) HandlerGetItems(w http.ResponseWriter, r *http.Request) {
	items, err := cfg.Queries.GetAllItems(context.Background())
	if err != nil {
		logger.Warn(err)
		return
	}

	respData, err := json.Marshal(items)
	if err != nil {
		logger.Warn(err)
		return
	}

	w.Write(respData)
}

func (cfg *ApiConfig) HandlerGetShoppingCart(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	userID, err := jwt.ValidateJWT(token, cfg.SecretJWT)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	shoppingCart, err := cfg.Queries.GetShoppingCart(context.Background(), userID)
	if err != nil {
		logger.Warn(err)
		return
	}

	respData, err := json.Marshal(shoppingCart)
	if err != nil {
		logger.Warn(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}

func (cfg *ApiConfig) HandlerGetInCart(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	userID, err := jwt.ValidateJWT(token, cfg.SecretJWT)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		logger.Warn(err)
		return
	}

	itemName, err := cfg.Queries.GetItemNameById(context.Background(), itemID)

	newItemInCart := GetItemParams{
		UserID: userID,
		ItemID: itemID,
		Name: itemName,
	}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&newItemInCart)
	if err != nil {
		logger.Warn(err)
		return
	}

	args := database.AddItemInCartParams{
		ItemID:   newItemInCart.ItemID,
		UserID:   newItemInCart.UserID,
		Quantity: int32(newItemInCart.Quantity),
		Cost:     int32(newItemInCart.Cost * newItemInCart.Quantity),
		ItemName: newItemInCart.Name,
	}

	err = cfg.Queries.AddItemInCart(context.Background(), args)
	if err != nil {
		logger.Warn(err)
		return
	}
}
