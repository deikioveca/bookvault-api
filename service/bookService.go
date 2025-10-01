package service

import (
	"BookVault-API/model"
	"errors"
	"gorm.io/gorm"
)

var (
	ErrBookNotFound 			= errors.New("book not found")
	ErrNoBooks					= errors.New("no books found")
)

type BookService interface{
	CreateBook(bookRequest *model.BookRequest) error

	GetByTitle(bookTitle string) (*model.BookResponse, error)

	GetBooks() ([]model.BookResponse, error)

	GetBooksByAuthor(author string) ([]model.BookResponse, error)

	UpdateStock(bookID uint) error
}

type bookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) BookService {
	return &bookService{db: db}
}


func (b *bookService) CreateBook(bookRequest *model.BookRequest) error {
	if bookRequest.Title == "" || bookRequest.Author == "" || bookRequest.Description == "" || bookRequest.Price == nil {
		return ErrEmptyFields
	}

	var book model.Book
	
	book.Title 			= bookRequest.Title
	book.Author 		= bookRequest.Author
	book.Description 	= bookRequest.Description
	book.Price 			= *bookRequest.Price
	book.InStock 		= true

	return b.db.Create(&book).Error
}


func (b *bookService) GetByTitle(bookTitle string) (*model.BookResponse, error) {
	var book model.Book

	if err := b.db.Where("LOWER(title) = LOWER(?)", bookTitle).First(&book).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	bookResponse := &model.BookResponse{
		Title: 			book.Title,
		Author: 		book.Author,
		Description: 	book.Description,
		Price: 			book.Price,
		InStock: 		book.InStock,
	}

	return bookResponse, nil
}


func (b *bookService) GetBooks() ([]model.BookResponse, error) {
	var books []model.Book
	var booksResponse []model.BookResponse

	if err := b.db.Find(&books).Error; err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, ErrNoBooks
	}

	for _, book := range books {
		bookResponse := model.BookResponse{
			Title: 			book.Title,
			Author: 		book.Author,
			Description: 	book.Description,
			Price: 			book.Price,
			InStock: 		book.InStock,
		}

		booksResponse = append(booksResponse, bookResponse)
	}

	return booksResponse, nil
}


func (b *bookService) GetBooksByAuthor(author string) ([]model.BookResponse, error) {
	var books []model.Book

	if err := b.db.Where("LOWER(author) = LOWER(?)", author).Find(&books).Error; err != nil {
		return nil, err
	}

	bookResponses := make([]model.BookResponse, 0, len(books))

	for _, book := range books {
		bookResponses = append(bookResponses, model.BookResponse{
			Title: 			book.Title,
			Author: 		book.Author,
			Description: 	book.Description,
			Price: 			book.Price,
			InStock: 		book.InStock,
		})
	}

	return bookResponses, nil
}


func (b *bookService) UpdateStock(bookID uint) error {
	var book model.Book

	if err := b.db.First(&book, bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return err
	}

	switch book.InStock {
	case true:
		book.InStock = false
	case false:
		book.InStock = true
	}

	return b.db.Save(&book).Error
}