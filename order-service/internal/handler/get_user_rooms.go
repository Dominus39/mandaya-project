package handler

import (
	"net/http"
	"order-service/config"
	"order-service/dto"
	"order-service/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// GetUserRooms godoc
// @Summary Get booked rooms for the authenticated user
// @Description Fetch all rooms currently booked by the authenticated user, including payment status.
// @Tags Rooms
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "List of booked rooms with payment status"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rooms/booked [get]
func GetUserRooms(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized access"})
	}

	claims, ok := user.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to parse user claims"})
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "User ID not found in claims"})
	}
	userID := int(userIDFloat)

	var bookings []models.Booking
	if err := config.DB.Preload("Room.Category").
		Where("user_id = ?", userID).
		Find(&bookings).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch booked rooms"})
	}

	if len(bookings) == 0 {
		return c.JSON(http.StatusOK, map[string]string{"message": "You have no booked room"})
	}

	var response []dto.GetUserRoomsResponse
	for _, booking := range bookings {
		var payment models.Booking
		isPaid := false
		if err := config.DB.Where("id = ?", booking.ID).First(&payment).Error; err == nil {
			isPaid = payment.IsPaid
		}

		response = append(response, dto.GetUserRoomsResponse{
			BookingID:  booking.ID,
			RoomID:     booking.Room.ID,
			RoomName:   booking.Room.Name,
			Category:   booking.Room.Category.Name,
			Price:      booking.Room.Category.Price,
			StartDate:  booking.StartDate,
			EndDate:    booking.EndDate,
			TotalPrice: booking.TotalPrice,
			IsPaid:     isPaid,
		})
	}

	return c.JSON(http.StatusOK, response)
}
