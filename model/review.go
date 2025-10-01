package model

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	Text		string
	UserID		uint
	User		User
	BookID		uint
	Book		Book
}

type ReviewRequest struct {
	Text		string	`json:"text"`
}

type ReviewResponse struct {
	Username		string	`json:"username"`
	Text			string	`json:"text"`
}

type UserReviewResponse struct {
	Username		string	`json:"username"`
	Title			string	`json:"title"`
	Author			string	`json:"author"`
	Text			string	`json:"text"`
}