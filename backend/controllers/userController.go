package controllers

import (
	"donation/helper"
	"donation/models"
	"donation/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userController struct {
	userService services.UserService
	authService services.AuthService
}

func NewUserController(userService services.UserService, authService services.AuthService) *userController {
	return &userController{userService, authService }
}

func (h *userController) RegisterUser(c *gin.Context) {
	var input helper.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrors(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ResponseJSON("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.ResponseJSON("Registered account failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		response := helper.ResponseJSON("Registered account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := helper.FormatUser(user, token)
	response := helper.ResponseJSON("Account hass ben registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userController) RegisterAdmin(c *gin.Context) {
	var input helper.RegisterAdminInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrors(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ResponseJSON("Register failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	registeredAdmin, err := h.userService.RegisterAdmin(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.ResponseJSON("Register failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(registeredAdmin.ID)
	if err != nil {
		response := helper.ResponseJSON("Register failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := helper.FormatUser(registeredAdmin, token)
	response := helper.ResponseJSON("Successfuly register", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userController) Login(c *gin.Context) {
	var input helper.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrors(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ResponseJSON("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.ResponseJSON("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.ResponseJSON("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := helper.FormatUser(loggedinUser, token)
	response := helper.ResponseJSON("Successfuly loggedin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userController) CheckEmailAvailability(c *gin.Context) {
	var input helper.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrors(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ResponseJSON("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server error"}
		response := helper.ResponseJSON("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	message := "Email has been registered"
	if isEmailAvailable {
		message = "Email is available"
	}

	response := helper.ResponseJSON(message, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userController) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		message := gin.H{"is_uploaded": false}
		response := helper.ResponseJSON("Failed to upload avatar image", http.StatusBadRequest, "error", message)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID
	path := fmt.Sprintf("../frontend/public/images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		message := gin.H{"is_uploaded": false}
		response := helper.ResponseJSON("Failed to upload avatar image", http.StatusBadRequest, "error", message)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, fmt.Sprintf("images/%d-%s", userID, file.Filename))
	if err != nil {
		message := gin.H{"is_uploaded": false}
		response := helper.ResponseJSON("Failed to upload avatar image", http.StatusBadRequest, "error", message)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	message := gin.H{"is_uploaded": true}
	response := helper.ResponseJSON("Avatar successfuly uploaded", http.StatusOK, "success", message)
	c.JSON(http.StatusOK, response)
}

func (h *userController) UpdateProfile(c *gin.Context) {
	var input helper.UpdateProfileInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrors(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ResponseJSON("Failed to update profile", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	ID := currentUser.ID 

	err = h.userService.UpdateProfile(input, ID)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ResponseJSON("Failed to update profile", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.ResponseJSON("Profile updated successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *userController) UpdatePassword(c *gin.Context) {
	var input helper.UpdatePasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
			errors := helper.ValidationErrors(err)
			errorMessage := gin.H{"errors": errors}
			response := helper.ResponseJSON("Failed to update password", http.StatusUnprocessableEntity, "error", errorMessage)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	ID := currentUser.ID

	if err := h.userService.UpdatePassword(input, ID); err != nil {
			errorMessage := gin.H{"errors": err.Error()}
			response := helper.ResponseJSON("Failed to update password", http.StatusUnprocessableEntity, "error", errorMessage)
			c.JSON(http.StatusUnprocessableEntity, response)
			return
	}

	response := helper.ResponseJSON("Password updated successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
