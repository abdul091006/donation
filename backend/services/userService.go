package services

import (
	"donation/helper"
	"donation/models"
	"donation/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(input helper.RegisterUserInput) (models.User, error)
	RegisterAdmin(input helper.RegisterAdminInput) (models.User, error)
	Login(input helper.LoginInput) (models.User, error)
	IsEmailAvailable(input helper.CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (models.User, error)
	GetUserByID(ID int) (models.User, error)
	UpdateProfile(input helper.UpdateProfileInput, ID int) error
	UpdatePassword(input helper.UpdatePasswordInput, ID int) error
}

type userService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *userService {
	return &userService{repository}
}

func (s *userService) RegisterUser(input helper.RegisterUserInput) (models.User, error) {
	user := models.User{}
	user.Name = input.Name
	user.Email = input.Email
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.Password = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (s *userService) RegisterAdmin(input helper.RegisterAdminInput) (models.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
			return user, err
	}

	if user.ID == 0 {
			return user, errors.New("no user found with that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}
	
	user.Role = "admin"

	updatedUser, err := s.repository.Update(user)
	if err != nil {
			return updatedUser, err
	}

	return updatedUser, nil
}



func (s *userService) Login(input helper.LoginInput) (models.User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *userService) IsEmailAvailable(input helper.CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *userService) SaveAvatar(ID int, fileLocation string) (models.User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	user.Avatar = fileLocation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *userService) GetUserByID(ID int) (models.User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("no user found with that ID")
	}

	return user, nil
}

func (s *userService) UpdateProfile(input helper.UpdateProfileInput, ID int) error {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return err
	}

	if input.Name != "" {
		user.Name = input.Name
	}

	if input.Email != "" {
		user.Email = input.Email
	}

	_, err = s.repository.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) UpdatePassword(input helper.UpdatePasswordInput, ID int) error {
	user, err := s.repository.FindByID(ID)
	if err != nil {
			return err
	}

	if input.NewPassword != input.ReEnterPassword {
		return errors.New("password tidak cocok")
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.MinCost)
	if err != nil {
			return err
	}

	user.Password = string(newPasswordHash)
	if _, err := s.repository.Update(user); err != nil {
			return err
	}

	return nil
}
