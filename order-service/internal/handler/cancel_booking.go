package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/config"
	"order-service/models"
	"order-service/utils"
	"strconv"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a user's booking by booking ID. Only the owner of the booking can cancel it.
// @Tags Rooms
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]string "Cancellation Successful"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 403 {object} map[string]string "Not authorized to cancel this booking"
// @Failure 404 {object} map[string]string "Booking not found"
// @Failure 500 {object} map[string]string "Cancellation failed"
// @Router /rooms/cancel/{id} [delete]
func CancelBooking(c echo.Context) error {
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

	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid booking ID"})
	}

	var booking models.Booking
	if err := config.DB.Preload("Room").First(&booking, bookingID).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Booking not found"})
	}

	if booking.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "You are not authorized to cancel this booking"})
	}

	var refundAmount float64
	if booking.IsPaid {
		paymentServiceURL := fmt.Sprintf("http://payment-service:8082/get_price/%d", bookingID)
		headers := map[string]string{"Content-Type": "application/json"}

		respBody, err := utils.RequestGET(paymentServiceURL, headers)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get booking price"})
		}

		var paymentResponse struct {
			Price float64 `json:"price"`
		}
		if err := json.Unmarshal(respBody, &paymentResponse); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to parse payment response"})
		}
		refundAmount = paymentResponse.Price

		userServiceURL := fmt.Sprintf("http://user-service:8080/update_balance/%d", userID)
		payload := map[string]float64{"amount": refundAmount}
		jsonData, _ := json.Marshal(payload)

		_, err = utils.RequestPOST(userServiceURL, headers, bytes.NewBuffer(jsonData))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user balance"})
		}

		booking.IsPaid = false
		if err := config.DB.Save(&booking).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update booking status"})
		}
	}

	if err := config.DB.Delete(&booking).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Cancellation failed"})
	}

	booking.Room.Stock += 1
	if err := config.DB.Save(&booking.Room).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update room stock"})
	}

	response := map[string]interface{}{
		"message": "Booking successfully cancelled",
	}

	if refundAmount > 0 {
		response["refund_amount"] = refundAmount
	}

	return c.JSON(http.StatusOK, response)
}
