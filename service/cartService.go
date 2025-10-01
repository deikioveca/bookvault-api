package service

import (
	"BookVault-API/model"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrCartNotFound 	= errors.New("cart not found")
	ErrCartBookNotFound = errors.New("book in cart not found")
)

type CartService interface {
	AddToCart(userID uint, bookID uint, quantity int) error

	ClearCart(userID uint) error

	RemoveFromCart(userID uint, bookID uint) error

	UpdateQuantity(userID uint, bookID uint, quantity int) error
	
	GetCart(cartID uint) (*model.CartResponse, error)
}

type cartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) CartService {
	return &cartService{db: db}
}


func (c *cartService) AddToCart(userID uint, bookID uint, quantity int) error {
	var cart model.Cart

	if err := c.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			cart = model.Cart{UserID: userID}
			if err := c.db.Create(&cart).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	var book model.Book

	if err := c.db.First(&book, bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return err
	}

	var cartBook model.CartBook

	if err := c.db.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).First(&cartBook).Error; err == nil {
		cartBook.Quantity += quantity
		return c.db.Save(&cartBook).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	cartBook = model.CartBook{
		CartID: cart.ID,
		BookID: bookID,
		Quantity: quantity,
	}

	return c.db.Create(&cartBook).Error
}


func (c *cartService) ClearCart(userID uint) error {
	var cart model.Cart

	if err := c.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartNotFound
		}
		return err
	}

	if err := c.db.Where("cart_id = ?", cart.ID).Delete(&model.CartBook{}).Error; err != nil {
		return err
	}

	return nil
}


func (c *cartService) RemoveFromCart(userID uint, bookID uint) error {
	var cart model.Cart

	if err := c.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartNotFound
		}
		return err
	}

	err := c.db.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).Delete(&model.CartBook{})
	if err.Error != nil {
		return err.Error
	}

	if err.RowsAffected == 0 {
		return ErrCartBookNotFound
	}
		
	return nil
}


func (c *cartService) UpdateQuantity(userID uint, bookID uint, quantity int) error {
	var cart model.Cart

	if err := c.db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartNotFound
		}
		return err
	}

	var cartBook model.CartBook

	if err := c.db.Where("cart_id = ? AND book_id = ?", cart.ID, bookID).First(&cartBook).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCartBookNotFound
		}
		return err
	}

	cartBook.Quantity = quantity

	return c.db.Save(&cartBook).Error
}


func (c *cartService) GetCart(cartID uint) (*model.CartResponse, error) {
	var cart model.Cart

	if err := c.db.Preload("Books.Book").First(&cart, cartID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	bookResponses := make([]model.BookResponse, 0, len(cart.Books))
	for _, b := range cart.Books {
		bookResponses = append(bookResponses, model.BookResponse{
			Title: 			b.Book.Title,
			Author: 		b.Book.Author,
			Description: 	b.Book.Description,
			Price: 			b.Book.Price,
			InStock: 		b.Book.InStock,
		})
	}

	return &model.CartResponse{ID: cart.ID, UserID: cart.UserID, Books: bookResponses}, nil
}