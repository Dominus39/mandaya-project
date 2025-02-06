package handler

import (
	"net/http"
	"order-service/config"
	"order-service/models"

	"github.com/labstack/echo/v4"
)

// Order-Service Update Payment Status Endpoint
func UpdatePaymentStatus(c echo.Context) error {
	// Get the booking ID from the URL parameter
	bookingID := c.Param("id")

	// Update the 'is_paid' status to true directly
	if err := config.DB.Model(&models.Booking{}).Where("id = ?", bookingID).Update("is_paid", true).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update payment status"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Payment status updated"})
}
