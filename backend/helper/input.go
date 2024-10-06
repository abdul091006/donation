package helper

import (
	"donation/models"
)

type UpdateProfileInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UpdatePasswordInput struct {
	NewPassword     string `json:"new_password" binding:"required"`
	ReEnterPassword     string `json:"re_enter_password" binding:"required"`
}

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterAdminInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type GetDonationDetailInput struct {
	ID int `uri:"id" binding:"required"`
}

type CreateDonationInput struct {
	Title            string                `form:"title" binding:"required"`
	ShortDescription string                `form:"short_description" binding:"required"`
	Description      string                `form:"description" binding:"required"`
	GoalAmount       int64                   `form:"goal_amount" binding:"required"`
	User             models.User
}

type UploadImageInput struct {
	ID           int `form:"id" binding:"required"`
}

type GetDonationTransactionInput struct {
	ID   int `uri:"id" binding:"required"`
	User models.User
}

type GetTransactionInput struct {
	ID   int `uri:"id" binding:"required"`
}

type CreateTransactionInput struct {
	Amount     int `json:"amount" bindiing:"required"`
	DonationID int `json:"donation_id" bindiing:"required"`
	User       models.User
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
