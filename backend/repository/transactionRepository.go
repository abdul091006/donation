package repository

import (
	"donation/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	GetByDonationID(donationID int) ([]models.Transaction, error)
	GetByUserID(userID int) ([]models.Transaction, error)
	GetByID(ID int) (models.Transaction, error)
	Save(transactionn models.Transaction) (models.Transaction, error)
	Update(transactionn models.Transaction) (models.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *transactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) GetByID(ID int) (models.Transaction, error) {
	var transaction models.Transaction

	err := r.db.Where("id = ?", ID).Preload("Donation").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) GetByDonationID(donationID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Where("donation_id = ?", donationID).Preload("User").Order("id DESC").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}



func (r *transactionRepository) GetByUserID(userID int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Where("user_id = ?", userID).Preload("Donation").Order("id DESC").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *transactionRepository) Save(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *transactionRepository) Update(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
