package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	UserID		uint
	User		User
	Status		string
	Address		string
	Total		float32
	Books		[]OrderBook
}


type OrderBook struct {
	gorm.Model
	OrderID		uint
	Order		Order
	BookID		uint
	Book		Book
	Quantity	uint
	Price		float32
}


type OrderResponse struct {
	ID			uint				`json:"id"`
	Status		string				`json:"status"`
	Total		float32				`json:"total"`
	Address	 	string	 			`json:"address"`
	CreatedAt	time.Time			`json:"created_at"`
	Books		[]OrderBookDetails	`json:"books"`	
}


type OrderBookDetails struct {
	Title    string  `json:"title"`
    Author   string  `json:"author"`
    Quantity int     `json:"quantity"`
    Price    float32 `json:"price"`
}