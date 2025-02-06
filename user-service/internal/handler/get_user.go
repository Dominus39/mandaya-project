package handler

import (
	"net/http"
	"user-service/config"
	"user-service/dto"
	"user-service/models"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	userID := c.Param("id")
	var user models.User

	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	response := dto.UserResponse{
		ID:      user.ID,
		Balance: user.Balance,
	}

	return c.JSON(http.StatusOK, response)
}
