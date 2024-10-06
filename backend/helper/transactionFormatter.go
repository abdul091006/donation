package helper

import (
	"donation/models"
	"time"
)

type DonationFormat struct {
	Title			string	`json:"title"`
	ImageURL	string	`json:"image_url"`
}

type DonationTransactionFormatter struct {
	ID        int    		`json:"id"`
	Name      string 		`json:"name"`
	Amount    int    		`json:"amount"`
	Status 		string		`json:"status"`
	CreatedAt time.Time	`json:"created_at"`
}

type TransactionByIDFormatter struct {
	ID        int    		`json:"id"`
	Amount    int    		`json:"amount"`
	Status		string		`json:"status"`
	CreatedAt time.Time	`json:"created_at"`
	Donation	DonationFormat	`json:"donation"`
}

func FormatDonationTransaction(transaction models.Transaction) DonationTransactionFormatter {
	formatter := DonationTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	return formatter
}

func FormatDonationTransactionByID(transaction models.Transaction) TransactionByIDFormatter {
	formatter := TransactionByIDFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	donationFormatter := DonationFormat{}
	donationFormatter.Title = transaction.Donation.Title
	donationFormatter.ImageURL = transaction.Donation.Image

	formatter.Donation = donationFormatter

	return formatter
}

func FormatDonationTransactions(transactions []models.Transaction) []DonationTransactionFormatter {
	var transactionsFormatter []DonationTransactionFormatter
	for _, transaction := range transactions{
		formatter := FormatDonationTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}

	return transactionsFormatter
}

type UserTransactionFormatter struct {
	ID				int							`json:"id"`
	Amount		int		 					`json:"amount"`
	Status		string					`json:"status"`
	CreatedAt	time.Time				`json:"created_at"`
	Donation	DonationFormat	`json:"donation"`
}


func FormatUserTransaction(transaction models.Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	donationFormatter := DonationFormat{}
	donationFormatter.Title = transaction.Donation.Title
	donationFormatter.ImageURL = transaction.Donation.Image

	formatter.Donation = donationFormatter

	return formatter
}

func FormatUserTransactions(transactions []models.Transaction) []UserTransactionFormatter {
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}

	var transactionFormatter []UserTransactionFormatter
	for _, transaction := range transactions{
		formatter := FormatUserTransaction(transaction)
		transactionFormatter = append(transactionFormatter, formatter)
	}

	return transactionFormatter
}

type TransactionFormatter struct {
	ID       		int    		`json:"id"`
	DonationID	int		 		`json:"donation_id"`
	UserID			int				`json:"user_id"`
	Amount    	int    		`json:"amount"`
	Status			string		`json:"status"`
	Code				string		`json:"code"`
	PaymentURL	string		`json:"payment_url"`
}

func FormatTransaction(transaction models.Transaction) TransactionFormatter {
	formatter := TransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.DonationID = transaction.DonationID
	formatter.UserID = transaction.UserID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.Code = transaction.Code
	formatter.PaymentURL =transaction.PaymentURL

	return formatter
}