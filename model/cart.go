package model

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID		uint		`gorm:"uniqueIndex"`
	User		User		`gorm:"constraint:OnDelete:CASCADE;"`
	Books		[]CartBook	`gorm:"constraint:OnDelete:CASCADE;"`
}


type CartBook struct {
	gorm.Model
	CartID		uint	`gorm:"index"`
	Cart		Cart	`gorm:"constraint:OnDelete:CASCADE;"`
	BookID		uint	`gorm:"index"`
	Book		Book	`gorm:"constraint:OnDelete:CASCADE;"`
	Quantity	int		
}


type CartResponse struct {
	ID		uint			`json:"id"`
	UserID	uint			`json:"user_id"`
	Books	[]BookResponse	`json:"books"`
}