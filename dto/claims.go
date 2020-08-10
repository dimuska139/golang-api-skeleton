package dto

import "github.com/dgrijalva/jwt-go"

type ClaimsDTO struct {
	jwt.StandardClaims
	User UserDTO `json:"user"`
}
