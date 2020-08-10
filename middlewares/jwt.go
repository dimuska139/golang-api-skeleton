package middlewares

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dimuska139/golang-api-skeleton/config"
	"github.com/dimuska139/golang-api-skeleton/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func parseToken(accessToken string, signingKey []byte) (*dto.UserDTO, error) {
	token, err := jwt.ParseWithClaims(accessToken, &dto.ClaimsDTO{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*dto.ClaimsDTO); ok && token.Valid {
		return &claims.User, nil
	}

	return nil, errors.New("Invalid token")
}

func JwtMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userDtoPtr, err := parseToken(parts[1], []byte(cfg.Jwt.Secret))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user", userDtoPtr)
	}
}
