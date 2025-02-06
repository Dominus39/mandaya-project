package dto

import "time"

type GetBookingResponse struct {
	BookingID  int       `json:"booking_id" gorm:"primaryKey;autoIncrement"`
	UserID     int       `json:"user_id" gorm:"not null"`
	RoomID     int       `json:"room_id" gorm:"not null"`
	RoomName   string    `json:"room_name" gorm:"-"`
	Category   string    `json:"category" gorm:"-"`
	Price      float64   `json:"price" gorm:"-"`
	StartDate  time.Time `json:"start_date" gorm:"not null"`
	EndDate    time.Time `json:"end_date" gorm:"not null"`
	TotalPrice float64   `json:"total_price" gorm:"not null"`
	IsPaid     bool      `json:"is_paid" gorm:"not null;default:false"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
}
