package helper

import (
	"donation/models"
)

type UserFormatter struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role string `json:"role"`
	Token string `json:"token"`
	ImageURL string `json:"image_url"`
}

func FormatUser(user models.User, token string) UserFormatter {
	formatter := UserFormatter{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Role: user.Role,
		Token: token,
		ImageURL: user.Avatar,
	}
	return formatter
}