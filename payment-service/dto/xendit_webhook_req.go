package dto

type XenditWebhookRequest struct {
	ID     string  `json:"id"`
	UserID int     `json:"user_id"`
	Amount float64 `json:"amount"`
	Status string  `json:"status"`
}
