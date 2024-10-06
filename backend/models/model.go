package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primarykey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Donation struct {
	ID               int       `json:"id" gorm:"primarykey"`
	UserID           int       `json:"user_id"`
	Title            string    `json:"title"`
	ShortDescription string    `json:"short_description"`
	Description      string    `json:"description"`
	BackerCount      int       `json:"backer_count"`
	GoalAmount       int64       `json:"goal_amount"`
	CurrentAmount    int       `json:"current_amount"`
	Slug             string    `json:"slug"`
	Image            string    `json:"image"`
	User             User      `json:"user"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Transaction struct {
	ID         int       `json:"id" gorm:"primarykey"`
	DonationID int       `json:"donation_id"`
	UserID     int       `json:"user_id"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Code       string    `json:"code"`
	PaymentURL string    `json:"payment_url"`
	User       User      `json:"user"`
	Donation   Donation  `json:"donation"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
