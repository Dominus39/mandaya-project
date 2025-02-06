package handler

import (
	"net/http"
	"order-service/config"
	"order-service/models"

	"github.com/labstack/echo/v4"
)

func UpdatePaymentStatus(c echo.Context) error {
	bookingID := c.Param("id")

	if err := config.DB.Model(&models.Booking{}).Where("id = ?", bookingID).Update("is_paid", true).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update payment status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Payment status updated"})
}
