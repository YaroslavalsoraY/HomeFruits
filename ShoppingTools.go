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
	Cost     int
}

func (cfg *ApiConfig) HandlerGetItems(w http.ResponseWriter, r *http.Request) {
	items, err := cfg.Queries.GetAllItems(context.Background())
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	respData, err := json.Marshal(items)
	if err != nil {
		http.Error(w, `{"error": "Problem with marshalling answer"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	w.Write(respData)
}

func (cfg *ApiConfig) HandlerGetShoppingCart(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	userID, err := jwt.ValidateJWT(token, cfg.SecretJWT)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	shoppingCart, err := cfg.Queries.GetShoppingCart(context.Background(), userID)
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	respData, err := json.Marshal(shoppingCart)
	if err != nil {
		http.Error(w, `{"error": "Problem with marshalling answer"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(respData)
}

func (cfg *ApiConfig) HandlerGetInCart(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	userID, err := jwt.ValidateJWT(token, cfg.SecretJWT)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, `{"error": "Problem with parsing request data"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	item, err := cfg.Queries.GetItemById(context.Background(), itemID)
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	newItemInCart := GetItemParams{
		UserID: userID,
		ItemID: itemID,
		Name:   item.Name,
		Cost:   int(item.Cost),
	}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&newItemInCart)
	if err != nil {
		http.Error(w, `{"error": "Problem with decoding json"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	if newItemInCart.Quantity > int(item.Quantity) {
		http.Error(w, `{"error": "Incorrect request data"}`, http.StatusBadRequest)
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
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	updateArgs := database.UpdateItemQuantityParams{
		Quantity: item.Quantity - int32(newItemInCart.Quantity),
		ID:       itemID,
	}

	err = cfg.Queries.UpdateItemQuantity(context.Background(), updateArgs)
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (cfg *ApiConfig) HandlerDeleteFromCart(w http.ResponseWriter, r *http.Request) {
	token, err := jwt.GetBearerToken(r.Header)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	userID, err := jwt.ValidateJWT(token, cfg.SecretJWT)
	if err != nil {
		http.Error(w, `{"error": "Unauthorized user"}`, http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	itemID, err := uuid.Parse(r.PathValue("itemID"))
	if err != nil {
		http.Error(w, `{"error": "Problem with parsing request data"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}
	
	args := database.DeleteFromCartParams{
		ItemID: itemID,
		UserID: userID,
	}
	deletedItem, err := cfg.Queries.DeleteFromCart(context.Background(), args)
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	realItem, err := cfg.Queries.GetItemById(context.Background(), deletedItem.ItemID)
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	err = cfg.Queries.UpdateItemQuantity(context.Background(), database.UpdateItemQuantityParams{
		Quantity: realItem.Quantity + deletedItem.Quantity,
		ID: deletedItem.ItemID,
	})
	if err != nil {
		http.Error(w, `{"error": "Problem with database query"}`, http.StatusInternalServerError)
		logger.Warn(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}