package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"payment-service/config"
	"payment-service/dto"
	"payment-service/models"
	"payment-service/utils"

	"github.com/labstack/echo/v4"
)

func XenditWebhook(c echo.Context) error {
	fmt.Println("Webhook endpoint hit!")

	var webhookData dto.XenditWebhookRequest

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to read request body"})
	}
	fmt.Printf("Raw Webhook Body: %s\n", string(body))

	if err := c.Bind(&webhookData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid webhook request"})
	}

	fmt.Printf("Webhook Data: %+v\n", webhookData)

	if webhookData.Status != "PAID" {
		fmt.Printf("Ignoring webhook, status: %s\n", webhookData.Status)
		return c.JSON(http.StatusOK, map[string]string{"message": "Webhook received but not paid"})
	}

	var paymentRecord models.PaymentForTopUp
	if err := config.DB.Where("invoice_id = ?", webhookData.ID).First(&paymentRecord).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Invoice not found"})
	}

	userServiceURL := fmt.Sprintf("http://localhost:8080/update_balance/%d", webhookData.UserID)
	payload := map[string]float64{
		"amount": webhookData.Amount,
	}
	headers := map[string]string{"Content-Type": "application/json"}

	jsonData, _ := json.Marshal(payload)
	_, err = utils.RequestPOST(userServiceURL, headers, bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user balance"})
	}

	paymentRecord.Status = "completed"
	config.DB.Save(&paymentRecord)

	fmt.Println("Top-up successful, user balance updated!")

	return c.JSON(http.StatusOK, map[string]string{"message": "Balance updated successfully"})
}
