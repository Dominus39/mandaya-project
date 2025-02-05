package models

type Category struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string  `json:"name" gorm:"not null;unique;size:100"`
	Description string  `json:"description" gorm:"size:255"`
	Price       float64 `json:"price" gorm:"not null"`
	Rooms       []Room  `json:"rooms" gorm:"foreignKey:CategoryID"`
}

type Room struct {
	ID         int      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" gorm:"not null;size:255"`
	CategoryID int      `json:"category_id" gorm:"not null"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
	Stock      int      `json:"stock" gorm:"not null"`
}
