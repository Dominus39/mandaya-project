package handler

import (
	"net/http"
	"payment-service/config"
	"payment-service/dto"
	"payment-service/models"

	"github.com/labstack/echo/v4"
)

func GetPrice(c echo.Context) error {
	bookingID := c.Param("id")
	var payment models.PaymentForBooking

	if err := config.DB.Where("id = ?", bookingID).First(&payment).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	response := dto.PriceResponse{
		Price: payment.Amount,
	}

	return c.JSON(http.StatusOK, response)
}
