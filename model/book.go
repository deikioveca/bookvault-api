package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string
	Author      string
	Description string
	Price       float32
	InStock     bool
	Reviews		[]Review
}


type BookRequest struct {
	Title		string		`json:"title"`
	Author		string		`json:"author"`
	Description	string		`json:"description"`
	Price		*float32	`json:"price"`
}


type BookResponse struct {
	Title		string	`json:"title"`
	Author		string	`json:"author"`
	Description	string	`json:"description"`
	Price		float32	`json:"price"`
	InStock		bool	`json:"in_stock"`
}