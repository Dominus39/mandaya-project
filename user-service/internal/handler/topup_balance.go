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
