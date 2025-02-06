package handler

import (
	"net/http"
	"user-service/config"
	"user-service/models"

	"github.com/labstack/echo/v4"
)

func UpdateUserBalance(c echo.Context) error {
	userID := c.Param("id")
	var payload struct {
		Amount float64 `json:"amount"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	user.Balance += payload.Amount
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update balance"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Balance updated successfully"})
}
