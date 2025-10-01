package services

import (
	"BookVault-API/model"
	"BookVault-API/service"
	"BookVault-API/tests/db"
	"testing"

	"gorm.io/gorm"
)

func initCartTestServices(t *testing.T) (*gorm.DB, service.CartService, service.BookService) {
	testDB := db.SetupTestDB(t)
	cartService := service.NewCartService(testDB)
	bookService := service.NewBookService(testDB)

	return testDB, cartService, bookService
}


func TestAddToCart(t *testing.T) {
	testDB, cartService, bookService := initCartTestServices(t)
	user := createTestUser(t, "testUser10001", "testUser@gmail.com")
	book := createTestBook(t, bookService, testDB ,"The Idiot", "Dostoevsky", 20.0)

	testCases := []struct{
		testName		string
		userID			uint
		bookID			uint
		quantity		int
		wantErr			error
		wantQuantity	int
	}{
		{"test new cart created and book added", user.ID, book.ID, 2, nil, 2},
		{"test add same book increases quantity", user.ID, book.ID, 1, nil, 3},
		{"test book not found", user.ID, 9999, 1, service.ErrBookNotFound, 0},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			err := cartService.AddToCart(test.userID, test.bookID, test.quantity)
			if err != test.wantErr {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
			if test.wantQuantity > 0 {
				var cartBook model.CartBook
				if err := testDB.First(&cartBook, "cart_id = ? AND book_id = ?", test.userID, test.bookID).Error; err != nil {
					t.Fatalf("failed to fetch cartBook: %v", err)
				}
				if cartBook.Quantity != test.wantQuantity {
					t.Errorf("expected quantity %d, got %d", test.wantQuantity, cartBook.Quantity)
				}
			}
		})
	}
}


func TestClearCart(t *testing.T) {
	testDB, cartService, bookService := initCartTestServices(t)

	user := createTestUser(t, "testUser5006", "testUser2006@gmail.com")
	book := createTestBook(t, bookService, testDB, "Notes From the Underground", "Fyodor Dostoevsky", 15.00)

	_ = cartService.AddToCart(user.ID, book.ID, 2)

	testCases := []struct{
		testName		string
		userID			uint
		wantErr			error
	}{
		{"test clear existing cart", user.ID, nil},
		{"test clear non-existent cart", 9999, service.ErrCartNotFound},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			err := cartService.ClearCart(test.userID)
			if err != test.wantErr {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
		})
	}
}


func TestRemoveFromCart(t *testing.T) {
	testDB, cartService, bookService := initCartTestServices(t)

	user := createTestUser(t, "testUser2000001", "testUser2000001@gmail.com")
	book := createTestBook(t, bookService, testDB, "The Gambler", "Фьодор Достоевски", 10.00)

	_ = cartService.AddToCart(user.ID, book.ID, 1)

	testCases := []struct{
		testName		string
		userID			uint
		bookID			uint
		wantErr			error
	}{
		{"test remove existing book", user.ID, book.ID, nil},
		{"test remove from non-existing cart", 999, book.ID, service.ErrCartNotFound},
		{"test remove non-existing book", user.ID, 999, service.ErrCartBookNotFound},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			err := cartService.RemoveFromCart(test.userID, test.bookID)
			if err != test.wantErr {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
		})
	}
}


func TestUpdateQuantity(t *testing.T) {
	testDB, cartService, bookService := initCartTestServices(t)

	user := createTestUser(t, "testUser2000001111", "testUser2000001111@gmail.com")
	book := createTestBook(t, bookService, testDB, "The Last Temptation", "Nikos Kazantsakis", 20.00)

	_ = cartService.AddToCart(user.ID, book.ID, 1)

	testCases := []struct{
		testName		string
		userID			uint
		bookID			uint
		newQuantity		int
		wantErr			error
		wantQuantity	int
	}{
		{"test update existing book quantity", user.ID, book.ID, 2, nil, 2},
		{"test update quantity in missing cart", 9999, book.ID, 3, service.ErrCartNotFound, 0},
		{"test missing book in existing cart", user.ID, 9999, 3, service.ErrCartBookNotFound, 0},
	}

	for _, test := range testCases{
		t.Run(test.testName, func(t *testing.T) {
			err := cartService.UpdateQuantity(test.userID, test.bookID, test.newQuantity)
			if err != test.wantErr {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
			if test.wantQuantity > 0 {
				var cartBook model.CartBook
				if err := testDB.First(&cartBook, "cart_id = ? AND book_id = ?", test.userID, test.bookID).Error; err != nil {
					t.Fatalf("failed to fetch cartBook: %v", err)
				}
				if cartBook.Quantity != test.wantQuantity {
					t.Errorf("expected quantity %d, got %d", test.wantQuantity, cartBook.Quantity)
				}
			}
		})
	}
}


func TestGetCart(t *testing.T){
	testDB, cartService, bookService := initCartTestServices(t)

	user := createTestUser(t, "testUser3400", "testUser3400@gmail.com")
	book := createTestBook(t, bookService, testDB, "1984", "George Orwell", 20.00)

	_ = cartService.AddToCart(user.ID, book.ID, 1)

	testCases := []struct{
		testName		string
		cartID			uint
		wantErr			error
		wantBook		string
	}{
		{"test get existing cart", 1, nil, "1984"},
		{"test get non-existing cart", 9999, service.ErrCartNotFound, ""},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			cart, err := cartService.GetCart(test.cartID)
			if err != test.wantErr {
				t.Errorf("expected error %v, got %v", test.wantErr, err)
			}
			if test.wantBook != "" && len(cart.Books) > 0 {
				if cart.Books[0].Title != test.wantBook {
					t.Errorf("expected %s, got %s", test.wantBook, cart.Books[0].Title)
				}
			}
		})
	}
}