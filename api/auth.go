package api

import (
	apiErrors "github.com/dimuska139/golang-api-skeleton/api_errors"
	"github.com/dimuska139/golang-api-skeleton/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthAPI struct {
	UsersService *services.UsersService
	AuthService  *services.AuthService
}

func NewAuthAPI(usersService *services.UsersService, authService *services.AuthService) *AuthAPI {
	return &AuthAPI{UsersService: usersService, AuthService: authService}
}

func (a *AuthAPI) Registration(c *gin.Context) {
	type RegistrationDTO struct {
		Email    string `json:"email" binding:"required"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var dto RegistrationDTO
	if err := c.BindJSON(&dto); err != nil { // TODO: Implement user friendly api_errors handling
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdUser, err := a.AuthService.Registration(dto.Email, dto.Name, dto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jwtPair, err := a.AuthService.GenerateJwtPair(*createdUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *jwtPair)
}

func (a *AuthAPI) Login(c *gin.Context) {
	type LoginDTO struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var dto LoginDTO
	if err := c.BindJSON(&dto); err != nil { // TODO: Implement user friendly api_errors handling
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authenticatedUser, err := a.AuthService.Login(dto.Email, dto.Password)
	if err != nil {
		if _, ok := err.(*apiErrors.NotFoundError); ok {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	jwtPair, err := a.AuthService.GenerateJwtPair(*authenticatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *jwtPair)
}

func (a *AuthAPI) RefreshTokens(c *gin.Context) {
	type refreshDTO struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	var dto refreshDTO
	if err := c.BindJSON(&dto); err != nil { // TODO: Implement user friendly api_errors handling
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userDto, err := a.UsersService.GetByToken(dto.RefreshToken)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Invalid refresh token",
		})
		return
	}

	jwtPair, err := a.AuthService.GenerateJwtPair(*userDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := a.AuthService.InvalidateOldToken(dto.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *jwtPair)
}
