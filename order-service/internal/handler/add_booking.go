package handler

import (
	"net/http"
	"order-service/config"
	"order-service/dto"
	"order-service/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// BookRoom godoc
// @Summary Book a room
// @Description Book a room for a given number of days and start date.
// @Tags Rooms
// @Accept json
// @Produce json
// @Param booking body dto.BookingRequest true "Booking Request"
// @Success 200 {object} dto.BookingResponse "Booking Successful"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 404 {object} map[string]string "Room not found"
// @Failure 500 {object} map[string]string "Booking failed"
// @Router /rooms/booking [post]
func BookRoom(c echo.Context) error {
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Unauthorized access"})
	}

	claims, ok := user.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to parse user claims"})
	}

	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "User ID not found in claims"})
	}
	userID := int(userIDFloat)

	var req dto.BookingRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request parameters"})
	}

	if req.RoomID == 0 || req.Days <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Room ID and number of days are required"})
	}

	var room models.Room
	if err := config.DB.Preload("Category").First(&room, req.RoomID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Room not found"})
	}

	if room.Stock <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Room is fully booked"})
	}

	// Calculate end date and total price
	endDate := req.StartDate.AddDate(0, 0, req.Days)
	totalPrice := float64(req.Days) * room.Category.Price

	newBooking := models.Booking{
		UserID:     userID,
		RoomID:     req.RoomID,
		StartDate:  req.StartDate,
		EndDate:    endDate,
		TotalPrice: totalPrice,
	}

	if err := config.DB.Create(&newBooking).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Booking failed"})
	}

	room.Stock -= 1
	if err := config.DB.Save(&room).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update room stock"})
	}

	response := dto.BookingResponse{
		Message:    "Order successful",
		RoomName:   room.Name,
		Category:   room.Category.Name,
		TotalPrice: totalPrice,
	}

	return c.JSON(http.StatusOK, response)
}
