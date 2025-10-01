package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username 	string
	Password	string
	Email		string
	Role		string
	Details		UserDetails	`gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Reviews		[]Review
}


type UserDetails struct {
	gorm.Model
	FullName	string
	PhoneNumber	string
	UserID		uint	`gorm:"uniqueIndex"`
}


type UserDetailsRequest struct {
	FullName 	string	`json:"full_name"`
	PhoneNumber	string	`json:"phone_number"`
}


type RegisterRequest struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
	Email		string	`json:"email"`
}


type LoginRequest struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}


type UserResponse struct {
	Username		string		`json:"username"`
	Email			string		`json:"email"`
	PhoneNumber		string		`json:"phone_number"`
	FullName		string		`json:"full_name"`
}