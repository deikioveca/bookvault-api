package services

import (
	"BookVault-API/model"
	"BookVault-API/service"
	"BookVault-API/tests/db"
	"errors"
	"testing"
	"gorm.io/gorm"
)

func float32Ptr(f float32) *float32 {
	return &f
}

func initBookTestServices(t *testing.T) (*gorm.DB, service.BookService){
	testDB := db.SetupTestDB(t)
	return testDB, service.NewBookService(testDB)
}

func createTestBook(t *testing.T, bookService service.BookService, testDB *gorm.DB, title, author string, price float32) model.Book {
	t.Helper()
	
	err := bookService.CreateBook(&model.BookRequest{
		Title: title,
		Author: author,
		Description: "Test Book",
		Price: &price,
	})
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}

	var book model.Book
	if err := testDB.Where("title = ?", title).First(&book).Error; err != nil {
		t.Fatalf("failed to fetch book: %v", err)
	}

	return book
}


func TestCreateBook(t *testing.T) {
	_, bookService := initBookTestServices(t)

	testCases := []struct{
		testName		string
		bookRequest		model.BookRequest
		wantErr			error
	}{
		{"test empty fields", model.BookRequest{Title: "", Author: "", Description: "test book", Price: nil}, service.ErrEmptyFields},
		{"test create book succeeds", model.BookRequest{Title: "The Idiot", Author: "Dostoevsky", Description: "The Idiot one of Dostoevsky 5 major books", Price: float32Ptr(20.00)}, nil},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			err := bookService.CreateBook(&test.bookRequest)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Errorf("expected %v, got %v", test.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("expected success, got %v", err)
			}
		})
	}
}


func TestGetByTitle(t *testing.T) {
	testDB, bookService := initBookTestServices(t)

	createTestBook(t, bookService, testDB, "The Idiot", "Dostoevsky", 20.00)

	testCases := []struct{
		testName		string
		bookTitle		string
		wantErr			error
		wantBook		*model.BookResponse
	}{
		{"test book not found", "Crime and Punishment", service.ErrBookNotFound, nil},
		{"test book found", "The Idiot", nil, &model.BookResponse{Title: "The Idiot", Author: "Dostoevsky", Description: "Test Book", Price: 20.00, InStock: true}},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			book, err := bookService.GetByTitle(test.bookTitle)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Errorf("expected %v, got %v", test.wantErr, err)
				}
				if book != nil {
					t.Errorf("expected nil book, got %v", err)
				}
				return
			}
			if err != nil {
				t.Errorf("expected success, got %v", err)
			}
			if book.Title != test.wantBook.Title {
				t.Errorf("expected %q, got %q", test.wantBook.Title, book.Title)
			}
		})
	}
}


func TestGetBooks(t *testing.T) {
	testDB, bookService := initBookTestServices(t)

	t.Run("test no books", func(t *testing.T) {
		books, err := bookService.GetBooks()
		if err == nil || err != service.ErrNoBooks {
			t.Errorf("expected ErrNoBooks, got %v", err)
		}
		if books != nil {
			t.Errorf("expected nil books, got %+v", books)
		}
	})

	t.Run("test books exist", func(t *testing.T) {
		createTestBook(t, bookService, testDB, "The Idiot", "Dostoevsky", 20.00)
		createTestBook(t, bookService, testDB, "Anna Karenina", "Tolstoy", 20.00)

		books, err := bookService.GetBooks()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(books) != 2 {
			t.Errorf("expected 2 books, got %d", len(books))
		}
	})
}


func TestGetBooksByAuthor(t *testing.T) {
	testDB, bookService := initBookTestServices(t)

	createTestBook(t, bookService, testDB, "The Idiot", "Dostoevsky", 20.00)
	createTestBook(t, bookService, testDB, "Crime and Punishment", "Dostoevsky", 20.00)
	createTestBook(t, bookService, testDB, "Ana Karenina", "Tolstoy", 20.00)
	createTestBook(t, bookService, testDB, "East of Eden", "Stainbeck", 20.00)

	testCases := []struct{
		testName		string
		bookAuthor		string
		wantCount		int
	}{
		{"test no books", "Balzac", 0},
		{"test 2 books with this author", "Dostoevsky", 2},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			books, err := bookService.GetBooksByAuthor(test.bookAuthor)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if len(books) != test.wantCount {
				t.Errorf("expected %d books, got %d", test.wantCount, len(books))
			}
		})
	}
}


func TestUpdateStock(t *testing.T) {
	testDB, bookService := initBookTestServices(t)

	book := createTestBook(t, bookService, testDB, "The Idiot", "Dostoevsky", 20.00)

	t.Run("test book not found", func(t *testing.T) {
		err := bookService.UpdateStock(9999)
		if err == nil || err != service.ErrBookNotFound {
			t.Errorf("expected ErrBookNotFound, got %v", err)
		}
	})

	t.Run("test stock update succeeds", func(t *testing.T) {
		initial := book.InStock

		if err := bookService.UpdateStock(book.ID); err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		var updatedBook model.Book
		if err := testDB.First(&updatedBook, book.ID).Error; err != nil {
			t.Fatalf("failed to fetch updated book: %v", err)
		}
		if updatedBook.InStock == initial {
			t.Errorf("expected stock change, got %v", updatedBook.InStock)
		}	
	})
}