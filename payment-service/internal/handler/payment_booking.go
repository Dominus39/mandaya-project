package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"payment-service/config"
	"payment-service/dto"
	"payment-service/models"

	"net/http"
	"payment-service/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// PayBooking godoc
// @Summary Pay for a booking
// @Description Pay the total price of a booking with the user's balance and mark it as paid.
// @Tags Payments
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {object} map[string]interface{} "Payment successful"
// @Failure 400 {object} map[string]string "Insufficient balance or invalid request"
// @Failure 404 {object} map[string]string "Booking not found"
// @Failure 500 {object} map[string]string "Payment failed"
// @Router /rooms/payment/{id} [post]
func PayBooking(c echo.Context) error {
	// Extract user claims from JWT
	user := c.Get("user")
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"message": "Unauthorized access"})
	}

	// Parse user claims as jwt.MapClaims
	claims, ok := user.(jwt.MapClaims)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed to parse user claims"})
	}

	// Extract user ID from claims
	userIDFloat, ok := claims["id"].(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "User ID not found in claims"})
	}
	userID := int(userIDFloat)

	// Parse the booking ID from the path
	bookingID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid booking ID"})
	}

	// Begin GORM transaction
	tx := config.DB.Begin()

	orderServiceURL := fmt.Sprintf("http://localhost:8081/get_booking/%d", bookingID)

	headers := map[string]string{
		"Authorization": c.Request().Header.Get("Authorization"),
	}

	res, err := utils.RequestGET(orderServiceURL, headers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "failed parsing response from server"})
	}

	var booking dto.BookingResponse
	if err := json.Unmarshal(res, &booking); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error while unmarshalling response"})
	}

	log.Printf("booking: %v", booking)
	log.Printf("booking.user_id: %d, user_id: %d", booking.UserID, userID)
	// Validate that the booking belongs to the user
	if booking.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "You are not authorized to pay for this booking"})
	}

	// Check if a payment already exists for this booking
	var payment models.PaymentForBooking
	paymentExists := true
	if err := tx.Where("booking_id = ?", bookingID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			paymentExists = false
		} else {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to query payment record"})
		}
	}

	// If payment exists and is already paid, return an error
	if paymentExists && booking.IsPaid {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Booking is already paid"})
	}

	// Find the user's account
	userServiceURL := fmt.Sprintf("http://localhost:8080/get_user/%d", userID)
	userRes, err := utils.RequestGET(userServiceURL, headers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch user details"})
	}

	var userAcc dto.UserResponse
	if err := json.Unmarshal(userRes, &userAcc); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error while unmarshalling user response"})
	}

	fmt.Printf("userAcc: %v\n", userAcc)
	fmt.Printf("userAcc.Balance: %v\n", userAcc.Balance)
	fmt.Printf("booking.TotalPrice: %v\n", booking.TotalPrice)
	// Check if the user has enough balance
	if userAcc.Balance < booking.TotalPrice {
		tx.Rollback()
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Insufficient balance"})
	}

	// Deduct the total price from the user's balance
	updateBalanceURL := fmt.Sprintf("http://localhost:8080/update_balance/%d", userID)
	updatePayload := map[string]float64{"amount": -booking.TotalPrice}

	// Convert map to JSON
	jsonData, err := json.Marshal(updatePayload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to encode request data"})
	}

	// Print the request payload and URL for debugging
	fmt.Printf("Request Payload: %s\n", jsonData)
	fmt.Printf("Request URL: %s\n", updateBalanceURL)

	bodyReader := bytes.NewBuffer(jsonData)
	headers["Content-Type"] = "application/json"

	_, err = utils.RequestPOST(updateBalanceURL, headers, bodyReader)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user balance"})
	}

	// Create or update the payment record
	if paymentExists {
		booking.IsPaid = true
		payment.CreatedAt = time.Now()
		if err := tx.Save(&payment).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update payment record"})
		}
	} else {
		booking.IsPaid = true
		payment = models.PaymentForBooking{
			BookingID: booking.BookingID,
			Amount:    booking.TotalPrice,
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&payment).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create payment record"})
		}
	}
	fmt.Printf("booking.BookingID: %v\n", booking.BookingID)
	// Save the booking updates
	updatePaymentStatusURL := fmt.Sprintf("http://localhost:8081/update_payment_status/%d", booking.BookingID)

	// Send the request to order-service to update the booking status
	_, err = utils.RequestGET(updatePaymentStatusURL, headers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update booking status"})
	}

	// Add a record to the user's history
	// userHistory := entity.UserHistory{
	// 	UserID:       userID,
	// 	Description:  "Payment for booking ID " + strconv.Itoa(bookingID),
	// 	ActivityType: "Payment",
	// 	ReferenceID:  bookingID,
	// }
	// if err := tx.Create(&userHistory).Error; err != nil {
	// 	tx.Rollback()
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to record user activity"})
	// }

	// Create the invoice via Xendit API
	// _, err = utils.CreateInvoice(booking, userAcc)
	// if err != nil {
	// 	tx.Rollback()
	// 	return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create invoice", "error": err.Error()})
	// }

	// Commit transaction
	tx.Commit()

	// Respond with success message
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Payment successful",
		"booking_id":  booking.BookingID,
		"total_price": booking.TotalPrice,
		"balance":     userAcc.Balance,
		"is_paid":     booking.IsPaid,
		"paid_at":     payment.CreatedAt,
	})
}
