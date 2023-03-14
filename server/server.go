package server

import (
	"context"
	"errors"
	"github.com/rs/cors"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pipeline1987/SVB/database"
	"github.com/pipeline1987/SVB/repositories"
	"github.com/pipeline1987/SVB/websocket"
)

type Config struct {
	PORT              string
	JWT_SECRET        string
	DB_HOST           string
	HASH_COST         string
	SIGN_EXPIRE_HOURS string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.PORT == "" {
		return nil, errors.New("PORT is required")
	}

	if config.JWT_SECRET == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	if config.DB_HOST == "" {
		return nil, errors.New("DB_HOST is required")
	}

	if config.HASH_COST == "" {
		return nil, errors.New("HASH_COST is required")
	}

	if config.SIGN_EXPIRE_HOURS == "" {
		return nil, errors.New("SIGN_EXPIRE_HOURS is required")
	}

	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}

	return broker, nil
}

func (b *Broker) Start(binder func(server Server, router *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)

	handler := cors.AllowAll().Handler(b.router)

	repo, err := database.NewPsqlRepository(b.config.DB_HOST)

	if err != nil {
		log.Fatal(err)
	}

	go b.hub.Run()
	repositories.SetRepository(repo)

	host := strings.Join([]string{
		"0.0.0.0:",
		b.config.PORT,
	}, "")

	log.Println("Server starting on port:", b.config.PORT)

	if err := http.ListenAndServe(host, handler); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
