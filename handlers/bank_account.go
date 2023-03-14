package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pipeline1987/SVB/models"
	"github.com/segmentio/ksuid"
	"net/http"

	"github.com/pipeline1987/SVB/middlewares"
	"github.com/pipeline1987/SVB/repositories"
	"github.com/pipeline1987/SVB/server"
)

type CreateBankAccountRequest struct {
	Name string `json:"name"`
}

type CreateBankAccountResponse struct {
	Id string `json:"id"`
}

type GetBankAccountResponse struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
	State   string  `json:"state"`
}

type UpdateBankAccountRequest struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

func CreateBankAccountHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(middlewares.ContextUserId)

		var request = CreateBankAccountRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		id, err := ksuid.NewRandom()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		var bankAccount = models.BankAccount{
			Id:      id.String(),
			UserId:  userId.(string),
			Name:    request.Name,
			Balance: 0,
			State:   "active",
		}

		savedBankAccount, repoErr := repositories.CreateBankAccount(r.Context(), &bankAccount)

		if repoErr != nil {
			http.Error(w, repoErr.Error(), http.StatusInternalServerError)

			return
		}

		var message = models.WebSocketMessage{
			Type:    "bank_account_created",
			Payload: savedBankAccount.Id,
		}

		s.Hub().Broadcast(message, nil)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateBankAccountResponse{
			Id: savedBankAccount.Id,
		})

	}
}

func GetBankAccountByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(middlewares.ContextUserId)
		params := mux.Vars(r)

		bankAccount, repoErr := repositories.GetBankAccountById(r.Context(), params["id"], userId.(string))

		if repoErr != nil {
			http.Error(w, repoErr.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetBankAccountResponse{
			Id:      bankAccount.Id,
			Name:    bankAccount.Name,
			Balance: bankAccount.Balance,
			State:   bankAccount.State,
		})
	}
}

func UpdateBankAccountByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(middlewares.ContextUserId)
		params := mux.Vars(r)

		var request = UpdateBankAccountRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		var bankAccount = models.BankAccount{
			Name:  request.Name,
			State: request.State,
		}

		updatedBankAccount, repoErr := repositories.UpdateBankAccountById(
			r.Context(),
			params["id"],
			userId.(string),
			&bankAccount,
		)

		if repoErr != nil {
			http.Error(w, repoErr.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetBankAccountResponse{
			Id:      updatedBankAccount.Id,
			Name:    updatedBankAccount.Name,
			Balance: updatedBankAccount.Balance,
			State:   updatedBankAccount.State,
		})
	}
}

func DeleteBankAccountByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(middlewares.ContextUserId)
		params := mux.Vars(r)

		repoErr := repositories.DeleteBankAccountById(
			r.Context(),
			params["id"],
			userId.(string),
		)

		if repoErr != nil {
			http.Error(w, repoErr.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

		return
	}
}

func GetAllBankAccountByUserIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(middlewares.ContextUserId)

		bankAccounts, repoErr := repositories.GetAllBankAccountsByUserId(r.Context(), userId.(string))

		if repoErr != nil {
			http.Error(w, repoErr.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(bankAccounts)
	}
}
