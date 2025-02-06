package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"payment-service/config"
	"payment-service/dto"
	"payment-service/models"
	"payment-service/utils"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateTopUpInvoice(c echo.Context) error {
	var req struct {
		UserID int     `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request parameters"})
	}

	userServiceURL := fmt.Sprintf("http://localhost:8080/get_user/%d", req.UserID)
	headers := map[string]string{
		"Authorization": c.Request().Header.Get("Authorization"),
		"Content-Type":  "application/json",
	}

	userRes, err := utils.RequestGET(userServiceURL, headers)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to fetch user details"})
	}

	var userAcc dto.UserDetailResponse
	if err := json.Unmarshal(userRes, &userAcc); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error while unmarshalling user response"})
	}

	invoice, err := utils.CreateInvoice(req.UserID, userAcc.Name, userAcc.Email, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create invoice"})
	}

	paymentRecord := models.PaymentForTopUp{
		UserID:    req.UserID,
		Amount:    req.Amount,
		InvoiceID: invoice.ID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&paymentRecord).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to store invoice details"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":     "Invoice created",
		"invoice_url": invoice.InvoiceURL,
	})
}
