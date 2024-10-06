package repository

import (
	"donation/models"

	"gorm.io/gorm"
)

type DonationRepository interface {
	FindAll() ([]models.Donation, error)
	FindByUserID(userID int) ([]models.Donation, error)
	FindByID(ID int) (models.Donation, error)
	Save(donation models.Donation) (models.Donation, error)
	Update(donation models.Donation) (models.Donation, error)
	Delete(ID int) error
}

type donationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) *donationRepository {
	return &donationRepository{db}
}

func (r *donationRepository) FindAll() ([]models.Donation, error) {
	var donations []models.Donation

	err := r.db.Find(&donations).Order("id DESC").Error
	if err != nil {
			return donations, err
	}

	return donations, nil
}


func (r *donationRepository) FindByUserID(userID int) ([]models.Donation, error) {
	var donations []models.Donation

	err := r.db.Where("user_id = ?", userID).Find(&donations).Error
	if err != nil {
		return donations, err
	}

	return donations, nil
}

func (r *donationRepository) FindByID(ID int) (models.Donation, error) {
	var donation models.Donation

	err := r.db.Where("id = ?", ID).Find(&donation).Error
	if err != nil {
		return donation, err
	}

	return donation, nil
}

func (r *donationRepository) Save(donation models.Donation) (models.Donation, error) {
	err := r.db.Create(&donation).Error
	if err != nil {
		return donation, err
	}

	return donation, nil
}

func (r *donationRepository) Update(donation models.Donation) (models.Donation, error) {
	err := r.db.Save(&donation).Error
	if err != nil {
		return donation, err
	}

	return donation, nil
}

func (r *donationRepository) Delete(ID int) error {
	var donation models.Donation

	if err := r.db.Where("id = ?", ID).Delete(&donation).Error; err != nil {
			return err 
	}


	return nil
}