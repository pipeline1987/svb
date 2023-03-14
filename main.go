package main

import (
	"context"
	"github.com/pipeline1987/SVB/middlewares"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pipeline1987/SVB/handlers"
	"github.com/pipeline1987/SVB/server"
)

func main() {
	envErr := godotenv.Load(".env")

	if envErr != nil {
		log.Fatalf("Error loading .env file%v\n", envErr)
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DB_HOST := os.Getenv("DB_HOST")
	HASH_COST := os.Getenv("HASH_COST")
	SIGN_EXPIRE_HOURS := os.Getenv("SIGN_EXPIRE_HOURS")

	s, serverErr := server.NewServer(context.Background(), &server.Config{
		PORT:              PORT,
		JWT_SECRET:        JWT_SECRET,
		DB_HOST:           DB_HOST,
		HASH_COST:         HASH_COST,
		SIGN_EXPIRE_HOURS: SIGN_EXPIRE_HOURS,
	})

	if serverErr != nil {
		log.Fatal("Error creating server instance &v\n", serverErr)
	}

	s.Start(BindRoutes)
}

func BindRoutes(s server.Server, r *mux.Router) {
	api := r.PathPrefix("/api").Subrouter()

	api.Use(middlewares.AuthMiddleware(s))

	api.HandleFunc("", handlers.HomeHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/users/sign-up", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/users/sign-in", handlers.SignInHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/users/me", handlers.GetUserHandler(s)).Methods(http.MethodGet)

	api.HandleFunc("/bank-accounts", handlers.CreateBankAccountHandler(s)).Methods(http.MethodPost)
	api.HandleFunc("/bank-accounts/{id}", handlers.GetBankAccountByIdHandler(s)).Methods(http.MethodGet)
	api.HandleFunc("/bank-accounts/{id}", handlers.UpdateBankAccountByIdHandler(s)).Methods(http.MethodPut)
	api.HandleFunc("/bank-accounts/{id}", handlers.DeleteBankAccountByIdHandler(s)).Methods(http.MethodDelete)
	api.HandleFunc("/bank-accounts", handlers.GetAllBankAccountByUserIdHandler(s)).Methods(http.MethodGet)

	api.HandleFunc("/ws", s.Hub().HandleWebSocket)
}
