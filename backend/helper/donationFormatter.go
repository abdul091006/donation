package helper

import (
	"donation/models"
)

type DonationFormatter struct {
	ID               int    `json:"id" gorm:"primarykey"`
	UserID           int    `json:"user_id"`
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
	Image         	 string `json:"image"`
	GoalAmount       int64    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

func FormatDonation(donation models.Donation) DonationFormatter {
	donationFormatter := DonationFormatter{}
	donationFormatter.ID = donation.ID
	donationFormatter.UserID = donation.UserID
	donationFormatter.Title = donation.Title
	donationFormatter.ShortDescription = donation.ShortDescription
	donationFormatter.GoalAmount = donation.GoalAmount
	donationFormatter.CurrentAmount = donation.CurrentAmount
	donationFormatter.Slug = donation.Slug
	donationFormatter.Image = donation.Image

	return donationFormatter
}

func FormatDonations(donations []models.Donation) []DonationFormatter {
	donationsFormatter := []DonationFormatter{}

	for _, donation := range donations {
		donationFormatter := FormatDonation(donation)
		donationsFormatter = append(donationsFormatter, donationFormatter)
	}

	return donationsFormatter
}

type DonationDetailFormatter struct {
	ID               int                   `json:"id" gorm:"primarykey"`
	UserID           int                   `json:"user_id"`
	Title            string                `json:"title"`
	ShortDescription string                `json:"short_description"`
	Description      string                `json:"description"`
	Image            string                `json:"image"`
	GoalAmount       int64                   `json:"goal_amount"`
	CurrentAmount    int                   `json:"current_amount"`
	Slug             string                `json:"slug"`
	User             DonationUserFormatter `json:"user"`
}

type DonationUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatDonationDetail(donation models.Donation) DonationDetailFormatter {
	donationDetailFormatter := DonationDetailFormatter{}
	donationDetailFormatter.ID = donation.ID
	donationDetailFormatter.UserID = donation.UserID
	donationDetailFormatter.Title = donation.Title
	donationDetailFormatter.ShortDescription = donation.ShortDescription
	donationDetailFormatter.Description = donation.Description
	donationDetailFormatter.GoalAmount = donation.GoalAmount
	donationDetailFormatter.CurrentAmount = donation.CurrentAmount
	donationDetailFormatter.Slug = donation.Slug
	donationDetailFormatter.Image = donation.Image

	user := donation.User
	donationUserFormatter := DonationUserFormatter{}
	donationUserFormatter.Name = user.Name
	donationUserFormatter.ImageURL = user.Avatar

	donationDetailFormatter.User = donationUserFormatter

	return donationDetailFormatter
}
