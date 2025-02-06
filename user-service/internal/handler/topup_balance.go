package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"user-service/dto"
	"user-service/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

// TopUpBalance handles the top-up balance request.
//
// @Summary Top up user balance
// @Description Allows authenticated users to top up their balance by creating an invoice via the payment service.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token"
// @Param request body dto.TopUpRequest true "Top-up request payload"
// @Success 200 {object} map[string]interface{} "Successful response with invoice details"
// @Failure 400 {object} map[string]string "Invalid request parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/topup [post]
func TopUpBalance(c echo.Context) error {
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

	var req dto.TopUpRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request parameters"})
	}

	if req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Amount must be greater than zero"})
	}

	paymentServiceURL := "http://localhost:8082/create_invoice"
	payload := map[string]interface{}{
		"user_id": userID,
		"amount":  req.Amount,
	}
	headers := map[string]string{
		"Authorization": c.Request().Header.Get("Authorization"),
		"Content-Type":  "application/json",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to encode request"})
	}

	respBody, err := utils.RequestPOST(paymentServiceURL, headers, bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create invoice"})
	}

	var resData map[string]interface{}
	if err := json.Unmarshal(respBody, &resData); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to parse response"})
	}

	return c.JSON(http.StatusOK, resData)
}
