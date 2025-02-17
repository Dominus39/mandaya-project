package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"payment-service/dto"
	"strconv"
)

func CreateInvoice(userID int, userName, userEmail string, amount float64) (*dto.InvoiceResponse, error) {
	apiKey := os.Getenv("XENDIT_API_SECRET")
	if apiKey == "" {
		return nil, fmt.Errorf("XENDIT_API_SECRET is not set")
	}
	apiUrl := os.Getenv("XENDIT_API_URL") + "/v2/invoices"

	if apiKey == "" || apiUrl == "" {
		return nil, errors.New("xendit API credentials not set")
	}

	bodyRequest := map[string]interface{}{
		"external_id":      "topup-" + strconv.Itoa(userID),
		"amount":           amount,
		"description":      fmt.Sprintf("Top-up balance for %s", userName),
		"invoice_duration": 86400, // 1-day invoice expiry
		"customer": map[string]interface{}{
			"name":  userName,
			"email": userEmail,
		},
		"currency":          "IDR",
		"should_send_email": true,
	}

	jsonData, err := json.Marshal(bodyRequest)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewBuffer(jsonData)

	encodedAPIKey := base64.StdEncoding.EncodeToString([]byte(apiKey + ":"))

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Basic " + encodedAPIKey, // Xendit requires Base64-encoded key
	}

	resBody, err := RequestPOST(apiUrl, headers, bodyReader)
	if err != nil {
		return nil, err
	}

	var resInvoice dto.InvoiceResponse
	if err := json.Unmarshal(resBody, &resInvoice); err != nil {
		return nil, err
	}

	return &resInvoice, nil
}
