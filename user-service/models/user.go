package models

type User struct {
	ID       int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string  `json:"name" gorm:"not null;size:100"`
	Email    string  `json:"email" gorm:"unique;not null;size:100"`
	Password string  `json:"-" gorm:"not null;size:100"`
	Balance  float64 `json:"balance" gorm:"not null;default:0"`
}
