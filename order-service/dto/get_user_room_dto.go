package dto

import "time"

type GetUserRoomsResponse struct {
	BookingID  int       `json:"booking_id"`
	RoomID     int       `json:"room_id"`
	RoomName   string    `json:"room_name"`
	Category   string    `json:"category"`
	Price      float64   `json:"price"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float64   `json:"total_price"`
	IsPaid     bool      `json:"is_paid"`
}
