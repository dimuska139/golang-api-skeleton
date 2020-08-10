package api

import (
	"github.com/dimuska139/golang-api-skeleton/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UsersAPI struct {
	UsersService *services.UsersService
}

func NewUsersAPI(usersService *services.UsersService) *UsersAPI {
	return &UsersAPI{UsersService: usersService}
}

func (a *UsersAPI) GetTotal(c *gin.Context) {
	total, e := a.UsersService.CountTotal()
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": total,
	})
}

func (a *UsersAPI) GetList(c *gin.Context) {
	users, e := a.UsersService.List()
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (a *UsersAPI) GetProfile(c *gin.Context) {
	userPtr, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, userPtr)
}
