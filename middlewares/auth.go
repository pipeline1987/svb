package middlewares

import (
	"context"
	"github.com/pipeline1987/SVB/repositories"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/pipeline1987/SVB/server"
)

var (
	NO_AUTH_NEEDED = []string{
		"api",
		"api/sign-up",
		"api/sign-in",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}

	return true
}

func AuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)

				return
			}

			tokenString := strings.TrimSpace(strings.Split(r.Header.Get("Authorization"), "Bearer")[1])

			parsedToken, jwtErr := jwt.ParseWithClaims(tokenString, &server.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWT_SECRET), nil
			})

			if jwtErr != nil {
				http.Error(w, jwtErr.Error(), http.StatusUnauthorized)

				return
			}

			if claims, ok := parsedToken.Claims.(*server.AppClaims); ok && parsedToken.Valid {
				if claims.UserId == "" {
					http.Error(w, "unauthorized", http.StatusUnauthorized)

					return
				}

				_, repoErr := repositories.ReadUser(r.Context(), claims.UserId)

				if repoErr != nil {
					http.Error(w, repoErr.Error(), http.StatusInternalServerError)
				}

				ctx := context.WithValue(r.Context(), ContextUserId, claims.UserId)

				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
