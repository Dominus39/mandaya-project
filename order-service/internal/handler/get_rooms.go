package handler

import (
	"net/http"
	"order-service/config"
	"order-service/dto"
	"order-service/models"

	"github.com/labstack/echo/v4"
)

// GetRooms godoc
// @Summary Get all available rooms
// @Description Get a list of all available rooms with name, category, price, and stock.
// @Tags Rooms
// @Accept json
// @Produce json
// @Success 200 {array} entity.RoomResponse "List of rooms"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /rooms [get]
func GetRooms(c echo.Context) error {
	var rooms []models.Room
	// Query all rooms with their category info and stock status
	err := config.DB.Preload("Category").Find(&rooms).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error fetching rooms"})
	}

	// Prepare a response structure with only required fields
	roomResponses := []dto.RoomResponse{}
	for _, room := range rooms {
		roomResponses = append(roomResponses, dto.RoomResponse{
			Name:     room.Name,
			Category: room.Category.Name,
			Price:    room.Category.Price,
			Stock:    room.Stock,
		})
	}

	return c.JSON(http.StatusOK, roomResponses)
}
