package services

import (
	"donation/helper"
	"donation/models"
	"donation/repository"
	"errors"
)

type TransactionService interface {
	GetTransactionByDonationID(input helper.GetDonationTransactionInput) ([]models.Transaction, error)
	GetTransactionByUserID(userID int) ([]models.Transaction, error)
	GetTransactionByID(input helper.GetTransactionInput) (models.Transaction, error)
	CreateTransaction(input helper.CreateTransactionInput) (models.Transaction, error)
}

type transactionService struct {
	repository         repository.TransactionRepository
	donationRepository repository.DonationRepository
	paymentService		 PaymentService
}

func NewTransactionService(repository repository.TransactionRepository, donationRepository repository.DonationRepository, paymentService PaymentService) *transactionService {
	return &transactionService{repository, donationRepository, paymentService}
}

func (s *transactionService) GetTransactionByDonationID(input helper.GetDonationTransactionInput) ([]models.Transaction, error) {
	donation, err := s.donationRepository.FindByID(input.ID)
	if err != nil {
		return []models.Transaction{}, err
	}

	if donation.UserID != input.User.ID {
		return []models.Transaction{}, errors.New("not an owner of the donation")
	}

	transactions, err := s.repository.GetByDonationID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *transactionService) GetTransactionByUserID(userID int) ([]models.Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *transactionService) GetTransactionByID(input helper.GetTransactionInput) (models.Transaction, error) {
	transaction, err := s.repository.GetByID(input.ID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (s *transactionService) CreateTransaction(input helper.CreateTransactionInput) (models.Transaction, error) {
	transaction := models.Transaction{}
	transaction.DonationID = input.DonationID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}

	paymentURL, err := s.paymentService.GetPaymentURL(newTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL
	newTransaction, err = s.repository.Update(newTransaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}