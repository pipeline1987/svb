package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pipeline1987/SVB/middlewares"
	"github.com/pipeline1987/SVB/models"
	"github.com/pipeline1987/SVB/repositories"
	"github.com/pipeline1987/SVB/server"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	Id string `json:"id"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
}

type GetUserResponse struct {
	Id       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpRequest{}
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

		defaultHash, _envErr := strconv.Atoi(s.Config().HASH_COST)

		if _envErr != nil {
			http.Error(w, _envErr.Error(), http.StatusInternalServerError)

			return
		}

		hashedPassword, cryptErr := bcrypt.GenerateFromPassword([]byte(request.Password), defaultHash)

		if cryptErr != nil {
			http.Error(w, cryptErr.Error(), http.StatusInternalServerError)

			return
		}

		var user = models.User{
			Id:       id.String(),
			Email:    request.Email,
			FullName: request.FullName,
			Password: string(hashedPassword),
		}

		savedUser, err := repositories.CreateUser(r.Context(), &user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignUpResponse{
			Id: savedUser.Id,
		})
	}
}

func SignInHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignInRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		user, repoErr := repositories.ReadUserByEmail(r.Context(), request.Email)

		if repoErr != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		if user == nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)

			return
		}

		if decryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); decryptErr != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)

			return
		}

		defaultSignExpireDays, _envErr := strconv.Atoi(s.Config().SIGN_EXPIRE_HOURS)

		if _envErr != nil {
			http.Error(w, _envErr.Error(), http.StatusInternalServerError)

			return
		}

		claims := server.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Duration(defaultSignExpireDays) * time.Hour * 24).Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, tokenErr := token.SignedString([]byte(s.Config().JWT_SECRET))

		if tokenErr != nil {
			http.Error(w, tokenErr.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignInResponse{
			AccessToken: tokenString,
		})
	}
}

func GetUserHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(middlewares.ContextUserId)

		user, repoErr := repositories.ReadUser(r.Context(), userId.(string))

		if repoErr != nil {
			http.Error(w, repoErr.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetUserResponse{
			Id:       user.Id,
			Email:    user.Email,
			FullName: user.FullName,
		})
	}
}
