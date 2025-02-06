package handler

import (
	"net/http"
	"order-service/config"
	"order-service/dto"
	"order-service/models"

	"github.com/labstack/echo/v4"
)

func GetBooking(c echo.Context) error {
	bookingID := c.Param("id")
	var booking models.Booking

	if err := config.DB.Preload("Room.Category").Where("id = ?", bookingID).First(&booking).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Booking not found"})
	}

	bookingResponse := dto.GetBookingResponse{
		BookingID:  booking.ID,
		UserID:     booking.UserID,
		RoomID:     booking.RoomID,
		RoomName:   booking.Room.Name,
		Category:   booking.Room.Category.Name,
		Price:      booking.Room.Category.Price,
		StartDate:  booking.StartDate,
		EndDate:    booking.EndDate,
		TotalPrice: booking.TotalPrice,
		IsPaid:     booking.IsPaid,
		CreatedAt:  booking.CreatedAt,
	}

	return c.JSON(http.StatusOK, bookingResponse)
}
