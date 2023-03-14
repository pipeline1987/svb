package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pipeline1987/SVB/server"
)

type HomeResponse struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

func HomeHandler(s server.Server) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(responseWriter).Encode(HomeResponse{
			Message: "Welcome to SVB",
			Status:  true,
		})

	}
}
