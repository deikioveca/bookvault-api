package service

import (
	"BookVault-API/model"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrEmptyCart 		= errors.New("cart is empty")
	ErrOrderNotFound 	= errors.New("order not found")
	ErrOrderCancel		= errors.New("order cannot be cancelled")
)

type OrderService interface {
	CreateOrder(userID uint, address string) (*model.OrderResponse, error)

	CancelOrder(orderID uint) error

	GetOrder(orderID uint) (*model.OrderResponse, error)

	GetUserOrders(userID uint) ([]model.OrderResponse, error)

	GetOrdersByStatus(status string) ([]model.OrderResponse, error)

	UpdateStatus(orderID uint, status string) error
}

type orderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) OrderService {
	return &orderService{db: db}
}


func (o *orderService) toOrderResponse(order *model.Order) *model.OrderResponse {
	var books []model.OrderBookDetails

	for _, book := range order.Books {
		books = append(books, model.OrderBookDetails{
			Title: 		book.Book.Title,
			Author: 	book.Book.Author,
			Quantity: 	int(book.Quantity),
			Price: 		book.Price,
		})
	}

	return &model.OrderResponse{
		ID: 		order.ID,
		Status: 	order.Status,
		Total: 		order.Total,
		Address: 	order.Address,
		CreatedAt: 	order.CreatedAt,
		Books: 		books,
	}
}


func (o *orderService) CreateOrder(userID uint, address string) (*model.OrderResponse, error) {
	var cart model.Cart
	
	if err := o.db.Preload("Books.Book").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEmptyCart
		}
		return nil, err
	}

	if len(cart.Books) == 0 {
		return nil, ErrEmptyCart
	}

	order := model.Order{
		UserID: 	userID,
		Status: 	"pending",
		Address: 	address,
	}

	var total float32
	var orderBooks []model.OrderBook

	for _, book := range cart.Books {
		linePrice := book.Book.Price * float32(book.Quantity)
		total += linePrice

		orderBooks = append(orderBooks, model.OrderBook{
			BookID: 	book.BookID,
			Quantity: 	uint(book.Quantity),
			Price: 		book.Book.Price,
		})
	}

	order.Total = total
	order.Books = orderBooks

	o.db.Create(&order)
	o.db.Where("cart_id = ?", cart.ID).Delete(&model.CartBook{})
	o.db.Preload("Books.Book").First(&order, order.ID)

	return o.toOrderResponse(&order), nil
}


func (o *orderService) CancelOrder(orderID uint) error {
	var order model.Order

	if err := o.db.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	if order.Status != "pending" && order.Status != "approved" {
		return ErrOrderCancel
	}

	order.Status = "cancelled"

	return o.db.Save(&order).Error
}


func (o *orderService) GetOrder(orderID uint) (*model.OrderResponse, error) {
	var order model.Order

	if err := o.db.Preload("Books.Book").First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return o.toOrderResponse(&order), nil
}


func (o *orderService) GetUserOrders(userID uint) ([]model.OrderResponse, error) {
	var orders []model.Order

	if err := o.db.Preload("Books.Book").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}

	var orderReponses []model.OrderResponse

	for _, order := range orders {
		orderReponses = append(orderReponses, *o.toOrderResponse(&order))
	}

	return orderReponses, nil
}


func (o *orderService) GetOrdersByStatus(status string) ([]model.OrderResponse, error) {
	var orders []model.Order

	if err := o.db.Preload("Books.Book").Where("LOWER(status) = LOWER(?)", status).Find(&orders).Error; err != nil {
		return nil, err
	}

	orderResponses := make([]model.OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, *o.toOrderResponse(&order))
	}

	return orderResponses, nil
}


func (o *orderService) UpdateStatus(orderID uint, status string) error {
	var order model.Order

	if err := o.db.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOrderNotFound
		}
		return err
	}

	order.Status = status

	return o.db.Save(&order).Error
}