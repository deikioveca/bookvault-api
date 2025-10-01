package service

import (
	"BookVault-API/model"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrEmptyReview 		= errors.New("review cannot be empty")
	ErrReviewNotFound 	= errors.New("review not found")
)

type ReviewService interface {
	AddReview(userID, bookID uint, reviewRequest model.ReviewRequest) error

	GetReviewsByBook(bookID uint) ([]model.ReviewResponse, error)

	GetReviewsByUser(userID	uint) ([]model.UserReviewResponse, error)

	UpdateReview(userID, bookID uint, reviewRequest model.ReviewRequest) error

	DeleteReviewByID(reviewID uint) error
}

type reviewService struct {
	db *gorm.DB
}

func NewReviewService(db *gorm.DB) ReviewService {
	return &reviewService{db: db}
}


func (r *reviewService) AddReview(userID, bookID uint, reviewRequest model.ReviewRequest) error {
	if reviewRequest.Text == "" {
		return ErrEmptyReview
	}

	var user model.User
	if err := r.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	var book model.Book
	if err := r.db.First(&book, bookID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookNotFound
		}
		return err
	}

	review := &model.Review{
		Text: 	reviewRequest.Text,
		UserID: userID,
		BookID: bookID,
	}
	return r.db.Create(review).Error
}


func (r *reviewService) GetReviewsByBook(bookID uint) ([]model.ReviewResponse, error) {
	var reviews []model.Review
	
	if err := r.db.Preload("User").Where("book_id = ?", bookID).Find(&reviews).Error; err != nil {
		return nil, err
	}

	reviewResponses := make([]model.ReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, model.ReviewResponse{
			Username: 	review.User.Username,
			Text: 		review.Text,
		})
	}

	return reviewResponses, nil
}


func (r *reviewService) GetReviewsByUser(userID	uint) ([]model.UserReviewResponse, error) {
	var reviews []model.Review

	if err := r.db.Preload("User").Preload("Book").Where("user_id = ?", userID).Find(&reviews).Error; err != nil {
		return nil, err
	}

	userReviewResponses := make([]model.UserReviewResponse, 0, len(reviews))
	for _, review := range reviews {
		userReviewResponses = append(userReviewResponses, model.UserReviewResponse{
			Username: 	review.User.Username,
			Title: 		review.Book.Title,
			Author: 	review.Book.Author,
			Text: 		review.Text,
		})
	}

	return userReviewResponses, nil
}


func (r *reviewService) UpdateReview(userID, bookID uint, reviewRequest model.ReviewRequest) error {
	var review model.Review

	if err := r.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&review).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		return err
	}

	review.Text = reviewRequest.Text

	return r.db.Save(&review).Error
}


func (r *reviewService) DeleteReviewByID(reviewID uint) error {
	var review model.Review

	if err := r.db.First(&review, reviewID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		return err
	}

	return r.db.Delete(&review).Error
}