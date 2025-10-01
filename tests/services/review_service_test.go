package services

import (
	"BookVault-API/model"
	"BookVault-API/service"
	"BookVault-API/tests/db"
	"errors"
	"testing"
	"gorm.io/gorm"
)

func initReviewTestServices(t *testing.T) (*gorm.DB, service.ReviewService, service.BookService, service.UserService) {
	t.Helper()

	testDB := db.SetupTestDB(t)
	reviewService 	:= service.NewReviewService(testDB)
	bookService 	:= service.NewBookService(testDB)
	userService 	:= service.NewUserService(testDB)

	return testDB, reviewService, bookService, userService
}

func createTestUserForReview(t *testing.T, testDB *gorm.DB, userService service.UserService, username, email string) model.User {
	t.Helper()

	if err := userService.Register(&model.RegisterRequest{
		Username: username,
		Password: "1234",
		Email:    email,
	}); err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	var user model.User
	if err := testDB.Where("username = ?", username).First(&user).Error; err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}

	return user
}

func createTestBookForReview(t *testing.T, testDB *gorm.DB, bookService service.BookService, title, author string, price float32) model.Book {
	t.Helper()

	err := bookService.CreateBook(&model.BookRequest{
		Title:       title,
		Author:      author,
		Description: "Test book",
		Price:       &price,
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


func TestAddReview(t *testing.T) {
	testDB, reviewService, bookService, userService := initReviewTestServices(t)

	user := createTestUserForReview(t, testDB, userService, "reviewuser", "reviewuser@gmail.com")
	book := createTestBookForReview(t, testDB, bookService, "Test Book", "Author", 12.50)

	testCases := []struct {
		name    			string
		userID  			uint
		bookID  			uint
		reviewRequest   	model.ReviewRequest
		wantErr 			error
	}{
		{"test empty text", user.ID, book.ID, model.ReviewRequest{Text: ""}, service.ErrEmptyReview},           
		{"test book not found", user.ID, 9999, model.ReviewRequest{Text: "Nice"}, service.ErrBookNotFound},
		{"test user not found", 9999, book.ID, model.ReviewRequest{Text: "Nice"}, service.ErrUserNotFound},
		{"test success", user.ID, book.ID, model.ReviewRequest{Text: "Loved it!"}, nil},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := reviewService.AddReview(test.userID, test.bookID, test.reviewRequest)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("expected error %v, got %v", test.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected success, got %v", err)
			}

			var review model.Review
			if err := testDB.Where("user_id = ? AND book_id = ?", test.userID, test.bookID).First(&review).Error; err != nil {
				t.Fatalf("expected to fetch review: %v", err)
			}

			if review.Text != test.reviewRequest.Text {
				t.Fatalf("expected text %q, got %q", test.reviewRequest.Text, review.Text)
			}
		})
	}
}


func TestGetReviewsByBook(t *testing.T) {
	testDB, reviewService, bookService, userService := initReviewTestServices(t)

	user := createTestUserForReview(t, testDB, userService, "reviewuser", "reviewuser@gmail.com")
	book := createTestBookForReview(t, testDB, bookService, "Test Book", "Author", 12.50) 
	
	err := reviewService.AddReview(user.ID, book.ID, model.ReviewRequest{Text: "Great book!"})
	if err != nil {
		t.Fatalf("failed to add review on book %d: %v", book.ID, err)
	}

	testCases := []struct{
		name		string
		bookID		uint
		wantCount	int
		wantErr		error
		wantText	[]string
	}{
		{"test valid book with reviews", book.ID, 1, nil, []string{"Great book!"}},
		{"test book with no reviews", 9999, 0, nil, []string{}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			reviews, err := reviewService.GetReviewsByBook(test.bookID)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("expected error %v, got %v", test.wantErr, err)
				}
				return 
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(reviews) != test.wantCount {
				t.Errorf("expected %d reviews, got %d", test.wantCount, len(reviews))
			}

			for _, want := range test.wantText {
				found := false
				for _, review := range reviews {
					if review.Text == want {
						found = true
						break
					}
				}

				if !found {
					t.Errorf("expected review text %q not found in %v", want, reviews)
				}
			}
		})
	}
}


func TestGetReviewsByUser(t *testing.T) {
	testDB, reviewService, bookService, userService := initReviewTestServices(t)

	user := createTestUserForReview(t, testDB, userService, "reviewuser3", "reviewuser3@gmail.com")
	book := createTestBookForReview(t, testDB, bookService, "Cousin Bette", "Balzac", 25.0)

	err := reviewService.AddReview(user.ID, book.ID, model.ReviewRequest{Text: "Amazing book!"})
	if err != nil {
		t.Fatalf("failed to add review on book: %v", err)
	}

	testCases := []struct{
		name		string
		userID		uint
		wantCount	int
		wantErr		error
		wantTexts	[]string
		wantTitles	[]string
	}{
		{"test valid user with reviews", user.ID, 1, nil, []string{"Amazing book!"}, []string{"Cousin Bette"}},
		{"test user with no reviews", 9999, 0, nil, []string{}, []string{}},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			reviews, err := reviewService.GetReviewsByUser(test.userID)
			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("expected error %v, got %v", test.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(reviews) != test.wantCount {
				t.Errorf("expected %d reviews, got %d", test.wantCount, len(reviews))
			}

			for _, want := range test.wantTexts {
				found := false
				for _, review := range reviews {
					if review.Text == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected review text %q not found in %+v", want, reviews)
				}
			}

			for _, want := range test.wantTitles {
				found := false
				for _, review := range reviews {
					if review.Title == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected review title %q not found in %+v", want, reviews)
				}
			}
		})
	}
}


func TestUpdateReview(t *testing.T) {
	testDB, reviewService, bookService, userService := initReviewTestServices(t)

    user := createTestUserForReview(t, testDB, userService,  "updateUser", "updateUser@gmail.com")
    book := createTestBookForReview(t, testDB, bookService, "Update Test Book", "Test Author", 30.00)

    err := reviewService.AddReview(user.ID, book.ID, model.ReviewRequest{Text: "Original review"})
    if err != nil {
        t.Fatalf("failed to add review on book: %v", err)
    }

	testCases := []struct{
		name			string
		userID			uint
		bookID			uint
		reviewRequest	model.ReviewRequest
		wantErr			error
	}{
		{"test valid update", user.ID, book.ID, model.ReviewRequest{Text: "Update review text"}, nil},
		{"test review not found", 9999, book.ID, model.ReviewRequest{Text: "Should not work"}, service.ErrReviewNotFound},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := reviewService.UpdateReview(test.userID, test.bookID, test.reviewRequest)

			if (err == nil) != (test.wantErr == nil) {
    			t.Fatalf("expected error %v, got %v", test.wantErr, err)
			}
			
			if err != nil && !errors.Is(err, test.wantErr) {
    			t.Fatalf("expected error %v, got %v", test.wantErr, err)
			}

			if err == nil {
				var review model.Review
				if dbErr := testDB.Where("user_id = ? AND book_id = ?", test.userID, test.bookID).First(&review).Error; dbErr != nil {
					t.Fatalf("failed to fetch review from db: %v", dbErr)
				}

				if review.Text != test.reviewRequest.Text {
					t.Errorf("expected review text %q, got %q", test.reviewRequest.Text, review.Text)
				}
			}
		})
	}
}


func TestDeleteReviewByID(t *testing.T) {
	testDB, reviewService, bookService, userService := initReviewTestServices(t)

    user := createTestUserForReview(t, testDB, userService,  "DeleteUser", "deleteUser@gmail.com")
    book := createTestBookForReview(t, testDB, bookService, "Delete Test Book", "Test Author", 30.00)

    err := reviewService.AddReview(user.ID, book.ID, model.ReviewRequest{Text: "Original review"})
    if err != nil {
        t.Fatalf("failed to add review on book: %v", err)
    }

	var review model.Review
	if err := testDB.Where("user_id = ? AND book_id = ?", user.ID, book.ID).First(&review).Error; err != nil {
		t.Fatalf("failed to fetch review: %v", err)
	}

	testCases := []struct{
		name		string
		reviewID	uint
		wantErr		error
	}{
		{"test delete existing review", review.ID, nil},
		{"test delete non-existing review", 9999, service.ErrReviewNotFound},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := reviewService.DeleteReviewByID(test.reviewID)

			if (err == nil) != (test.wantErr == nil) {
				t.Fatalf("expected error %v, got %v", test.wantErr, err)
			}

			if err != nil && !errors.Is(err, test.wantErr) {
				t.Fatalf("expected error %v, got %v", test.wantErr, err)
			}

			if test.wantErr == nil {
				var check model.Review
				if err := testDB.First(&check, test.reviewID).Error; err == nil {
					t.Errorf("expected review to be deleted, but it still exists")
				}
			}
		})
	}
}