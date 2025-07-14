package main

import (
	"HomeFruits/internal/database"
	"HomeFruits/logger"
	"database/sql"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ApiConfig struct {
	Queries    *database.Queries
	SecretJWT  string
	AdminEmail string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.HaltOnErr(err)
	}

	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		logger.HaltOnErr(err)
	}

	secretJWT := os.Getenv("SECRET_JWT")

	adminEmail := os.Getenv("ADMIN_EMAIL")

	config := ApiConfig{
		Queries:    database.New(db),
		SecretJWT:  secretJWT,
		AdminEmail: adminEmail,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/items", config.HandlerGetItems)
	mux.HandleFunc("GET /api/shopping_cart", config.HandlerGetShoppingCart)
	mux.HandleFunc("GET /api/delete/{itemID}", config.HandlerDeleteFromCart)

	mux.HandleFunc("POST /api/reg", config.HandlerRegUser)
	mux.HandleFunc("POST /api/login", config.HandlerLogin)
	mux.HandleFunc("POST /api/item/{itemID}", config.HandlerGetInCart)
	mux.HandleFunc("POST /api/refresh", config.HandlerRefresh)

	mux.HandleFunc("POST /admin/item", config.HandlerInsertItem)
	mux.HandleFunc("GET /admin/revoke/{tokenID}", config.HandlerRevokeToken)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Info("Starting server...")
	err = server.ListenAndServe()
	if err != nil {
		logger.HaltOnErr(err)
	}
}
