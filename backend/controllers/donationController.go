package controllers

import (
	"donation/helper"
	"donation/models"
	"donation/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type donationController struct {
	service services.DonationService
}

func NewDonationController(service services.DonationService) *donationController {
	return &donationController{service}
}

func (h *donationController) GetDonations(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	donations, err := h.service.GetDonations(userID)
	if err != nil {
		response := helper.ResponseJSON("Failed to get donations", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseJSON("List of donations", http.StatusOK, "success", helper.FormatDonations(donations))
	c.JSON(http.StatusOK, response)
}

func (h *donationController) GetDonation(c *gin.Context) {
	var input helper.GetDonationDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ResponseJSON("Failed to get detail of donation", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	donation, err := h.service.GetDonationByID(input)
	if err != nil {
		response := helper.ResponseJSON("Failed to get detail of donation", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseJSON("Donation Detail", http.StatusOK, "success", helper.FormatDonationDetail(donation))
	c.JSON(http.StatusOK, response)
}

func (h *donationController) CreateDonation(c *gin.Context) {
	var input helper.CreateDonationInput
	if err := c.ShouldBind(&input); err != nil {
		errors := helper.ValidationErrors(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ResponseJSON("Failed to parse JSON data", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		message := gin.H{"is_uploaded": false}
		response := helper.ResponseJSON("Failed to upload donation image", http.StatusBadRequest, "error", message)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID := c.MustGet("currentUser").(models.User).ID
	path := fmt.Sprintf("../frontend/public/images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		message := gin.H{"is_uploaded": false}
		response := helper.ResponseJSON("Failed to upload donation image", http.StatusBadRequest, "error", message)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	input.User = c.MustGet("currentUser").(models.User)

	newDonation, err := h.service.CreateDonation(input, fmt.Sprintf("images/%d-%s", userID, file.Filename))
	if err != nil {
		response := helper.ResponseJSON("Failed to create donation", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.ResponseJSON("Success to create donation", http.StatusOK, "success", helper.FormatDonation(newDonation))
	c.JSON(http.StatusOK, response)
}

func (h *donationController) GetDonationsByUserID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		response := helper.ResponseJSON("Invalid user ID", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	donations, err := h.service.GetDonationsByUserID(userID)
	if err != nil {
		response := helper.ResponseJSON("Failed to get user donations", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.ResponseJSON("User donations retrieved successfully", http.StatusOK, "success", helper.FormatDonations(donations))
	c.JSON(http.StatusOK, response)
}

func (h *donationController) UpdateDonation(c *gin.Context) {
	var inputID helper.GetDonationDetailInput

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.ResponseJSON("Failed to update donation", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData helper.CreateDonationInput

	err = c.ShouldBind(&inputData)
	if err != nil {
		errors := helper.ValidationErrors(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ResponseJSON("Failed to update donation", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		message := gin.H{"is_uploaded": false}
		response := helper.ResponseJSON("Failed to upload donation image", http.StatusBadRequest, "error", message)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	userID := c.MustGet("currentUser").(models.User).ID
	path := fmt.Sprintf("../frontend/public/images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		message := gin.H{"is_uploaded": false}
		response := helper.ResponseJSON("Failed to upload donation image", http.StatusBadRequest, "error", message)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	inputData.User = currentUser

	updatedDonation, err := h.service.UpdateDonation(inputID, inputData, fmt.Sprintf("images/%d-%s", userID, file.Filename))
	if err != nil {
		response := helper.ResponseJSON("Failed to update donation", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseJSON("Success to update donation", http.StatusOK, "success", helper.FormatDonation(updatedDonation))
	c.JSON(http.StatusOK, response)
}

func (h *donationController) DeleteDonation(c *gin.Context) {
	donationID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response := helper.ResponseJSON("Invalid donation ID", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	err = h.service.DeleteDonation(donationID)
	if err != nil {
		response := helper.ResponseJSON("Failed to delete donation", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.ResponseJSON("Donation deleted successfully", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
