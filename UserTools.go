package main

import (
	"HomeFruits/internal/database"
	"HomeFruits/internal/hashFunc"
	"HomeFruits/internal/jwt"
	"HomeFruits/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Token        string `json:"JWT,omitempty"`
	RefreshToken string
}

type Token struct {
	Token string `json:"refresh_token"`
}

const EXPIRESEIN = 15

func (cfg *ApiConfig) HandlerRegUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	newUser := User{}
	err := decoder.Decode(&newUser)
	if err != nil {
		logger.Warn(err, "problem with registration")
		return
	}

	hashedPassword, err := hashfunc.HashingPassword(newUser.Password)
	if err != nil {
		logger.Warn(err, "problem with hashing")
		return
	}

	arg := database.CreateNewUserParams{
		Email:          newUser.Email,
		HashedPassword: hashedPassword,
	}

	createdUser, err := cfg.Queries.CreateNewUser(context.Background(), arg)
	if err != nil {
		logger.Warn(err)
		return
	}

	token, err := jwt.MakeJWT(createdUser.ID, cfg.SecretJWT, EXPIRESEIN*time.Minute)
	if err != nil {
		logger.Warn(err)
		return
	}

	refreshToken := jwt.MakeRefreshToken()
	args := database.InsertNewRefreshTokenParams{
		Token:     refreshToken,
		UserID:    createdUser.ID,
		ExpiresAt: time.Now().Add(EXPIRESEIN * time.Hour * 24),
	}
	err = cfg.Queries.InsertNewRefreshToken(context.Background(), args)
	if err != nil {
		logger.Warn(err)
		return
	}
	newUser.RefreshToken = refreshToken

	logger.Info("New user created!")

	newUser.Token = token
	newUser.Password = "*********"
	respData, err := json.Marshal(newUser)

	w.WriteHeader(http.StatusCreated)
	w.Write(respData)
}

func (cfg *ApiConfig) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userData := User{}
	err := decoder.Decode(&userData)
	if err != nil {
		logger.Warn(err)
		return
	}

	realPasswordAndId, err := cfg.Queries.GetUserPassword(context.Background(), userData.Email)
	if err != nil {
		logger.Warn(err)
		return
	}

	if !hashfunc.HashCompareWithPassw(userData.Password, realPasswordAndId.HashedPassword) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	logger.Info(fmt.Sprintf("User: %s logged in", userData.Email))

	token, err := jwt.MakeJWT(realPasswordAndId.ID, cfg.SecretJWT, EXPIRESEIN*time.Minute)
	if err != nil {
		logger.Warn(err)
		return
	}

	refreshToken := jwt.MakeRefreshToken()
	args := database.InsertNewRefreshTokenParams{
		Token:     refreshToken,
		UserID:    realPasswordAndId.ID,
		ExpiresAt: time.Now().Add(EXPIRESEIN * time.Hour * 24),
	}
	err = cfg.Queries.InsertNewRefreshToken(context.Background(), args)
	if err != nil {
		logger.Warn(err)
		return
	}

	userData.Password = "*********"
	userData.Token = token
	userData.RefreshToken = refreshToken

	respData, err := json.Marshal(userData)

	w.Write(respData)
}

func (cfg *ApiConfig) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	refreshToken := Token{}
	err := decoder.Decode(&refreshToken)
	if err != nil {
		logger.Warn(err)
		return
	}

	tokenInfo, err := cfg.Queries.GetRefreshToken(context.Background(), refreshToken.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		logger.Warn(err)
		return
	}

	if tokenInfo.RevokedAt.Valid || time.Now().After(tokenInfo.ExpiresAt) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newAccessToken, err := jwt.MakeJWT(tokenInfo.UserID, cfg.SecretJWT, EXPIRESEIN * time.Minute)
	if err != nil {
		logger.Warn(err)
		return
	}

	respData, err := json.Marshal(newAccessToken)
	if err != nil {
		logger.Warn(err)
		return
	}

	w.Write(respData)
}
