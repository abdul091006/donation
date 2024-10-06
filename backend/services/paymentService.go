package services

import (
	"donation/helper"
	"donation/models"
	"donation/repository"
	"fmt"
	"strconv"

	"github.com/veritrans/go-midtrans"
)

type PaymentService interface {
	GetPaymentURL(transaction models.Transaction, user models.User) (string, error)
	ProcessPayment(input helper.TransactionNotificationInput) error
}

type paymentService struct {
	transactionRepository repository.TransactionRepository
	donationRepository    repository.DonationRepository
}

func NewPaymentServicce(transactionRepository repository.TransactionRepository, donationRepository repository.DonationRepository) *paymentService {
	return &paymentService{transactionRepository, donationRepository}
}

func (s *paymentService) GetPaymentURL(transaction models.Transaction, user models.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "RAHASIA"
	midclient.ClientKey = "RAHASIA"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			FName: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}

func (s *paymentService) ProcessPayment(input helper.TransactionNotificationInput) error {
	transactionID, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.transactionRepository.GetByID(transactionID)
	if err != nil {
		return err
	}

	fmt.Println(input.PaymentType)
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.transactionRepository.Update(transaction)
	if err != nil {
		return err
	}

	donation, err := s.donationRepository.FindByID(updatedTransaction.DonationID)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		donation.BackerCount = donation.BackerCount + 1
		donation.CurrentAmount = donation.CurrentAmount + updatedTransaction.Amount

		_, err := s.donationRepository.Update(donation)
		if err != nil {
			return nil
		}
	}

	return nil
}
