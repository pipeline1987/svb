package server

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	UserId string
	jwt.StandardClaims
}
