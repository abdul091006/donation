package services

import (
	"donation/helper"
	"donation/models"
	"donation/repository"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
)

type DonationService interface {
	GetDonations(userID int) ([]models.Donation, error)
	GetDonationByID(input helper.GetDonationDetailInput) (models.Donation, error)
	GetDonationsByUserID(userID int) ([]models.Donation, error)
	CreateDonation(input helper.CreateDonationInput, fileLocation string) (models.Donation, error)
	UpdateDonation(inputID helper.GetDonationDetailInput, inputData helper.CreateDonationInput, fileLocation string) (models.Donation, error)
	DeleteDonation(ID int) error
}

type donationService struct {
	repository repository.DonationRepository
}

func NewDonationService(repository repository.DonationRepository) *donationService {
	return &donationService{repository}
}

func (s *donationService) GetDonations(userID int) ([]models.Donation, error) {
	if userID != 0 {
		donations, err := s.repository.FindByUserID(userID)
		if err != nil {
			return donations, err
		}

		return donations, nil
	}

	donations, err := s.repository.FindAll()
	if err != nil {
		return donations, err
	}

	return donations, nil
}

func (s *donationService) GetDonationByID(input helper.GetDonationDetailInput) (models.Donation, error) {
	donation, err := s.repository.FindByID(input.ID)

	if err != nil {
		return donation, err
	}

	return donation, nil
}

func (s *donationService) CreateDonation(input helper.CreateDonationInput, fileLocation string) (models.Donation, error) {
	donation := models.Donation{}
	donation.Title = input.Title
	donation.ShortDescription = input.ShortDescription
	donation.Description = input.Description
	donation.GoalAmount = input.GoalAmount
	donation.Image = fileLocation
	donation.UserID = input.User.ID

	strSlug := fmt.Sprintf("%s %d", input.Title, input.User.ID)
	donation.Slug = slug.Make(strSlug)

	newDonation, err := s.repository.Save(donation)
	if err != nil {
		return newDonation, err
	}

	return newDonation, nil
}

func (s *donationService) GetDonationsByUserID(userID int) ([]models.Donation, error) {
	donations, err := s.repository.FindByUserID(userID)
	if err != nil {
		return donations, err
	}

	return donations, nil
}

func (s *donationService) UpdateDonation(inputID helper.GetDonationDetailInput, inputData helper.CreateDonationInput, fileLocation string) (models.Donation, error) {
	donation, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return donation, err
	}

	if donation.UserID != inputData.User.ID {
		return donation, errors.New("not an owner of the donation")
	}

	donation.Title = inputData.Title
	donation.ShortDescription = inputData.ShortDescription
	donation.Description = inputData.Description
	donation.Image = fileLocation
	donation.GoalAmount = inputData.GoalAmount

	updatedDonation, err := s.repository.Update(donation)
	if err != nil {
		return updatedDonation, err
	}

	return updatedDonation, nil
}

func (s *donationService) DeleteDonation(ID int) error {
	_, err := s.repository.FindByID(ID)
	if err != nil {
		return err
	}

	err = s.repository.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}