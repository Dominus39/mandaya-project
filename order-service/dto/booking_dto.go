package dto

import "time"

type BookingRequest struct {
	RoomID    int       `json:"room_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	Days      int       `json:"days" validate:"required"`
}

type BookingResponse struct {
	Message    string  `json:"message"`
	RoomName   string  `json:"room_name"`
	Category   string  `json:"category"`
	TotalPrice float64 `json:"total_price"`
}
