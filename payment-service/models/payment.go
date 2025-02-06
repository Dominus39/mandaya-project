package models

import "time"

type PaymentForBooking struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	BookingID int       `json:"booking_id" gorm:"not null;unique"`
	Amount    float64   `json:"amount" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type PaymentForTopUp struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int       `json:"user_id" gorm:"not null"`
	Amount    float64   `json:"amount" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
