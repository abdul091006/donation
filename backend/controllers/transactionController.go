package controllers

import (
	"donation/helper"
	"donation/models"
	"donation/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionsController struct {
	service services.TransactionService
	paymentService services.PaymentService
}

func NewTransactionController(service services.TransactionService, paymentService services.PaymentService) *transactionsController {
	return &transactionsController{service, paymentService}
}

func (h *transactionsController) GetDonationTransations(c *gin.Context) {
	var input helper.GetDonationTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ResponseJSON("Failed to get donation transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionByDonationID(input)
	if err != nil {
		response := helper.ResponseJSON("Failed to get donation transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseJSON("List of transactions", http.StatusOK, "success", helper.FormatDonationTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionsController) GetTransactionByID(c *gin.Context) {
	var input helper.GetTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ResponseJSON("Failed to get donation transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	transactions, err := h.service.GetTransactionByID(input)
	if err != nil {
		response := helper.ResponseJSON("Failed to get donation transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseJSON("List of transactions", http.StatusOK, "success", helper.FormatDonationTransactionByID(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionsController) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionByUserID(userID)
	if err != nil {
		response := helper.ResponseJSON("Failed to get user's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseJSON("User's transactions", http.StatusOK, "success", helper.FormatUserTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionsController) CreateTransaction(c *gin.Context) {
	var input helper.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.ValidationErrors(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ResponseJSON("Failed to create transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.ResponseJSON("Failed to create transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ResponseJSON("Success to create transaction", http.StatusOK, "success", helper.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)
}

func (h *transactionsController) GetNotification(c *gin.Context) {
	var input helper.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.ResponseJSON("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	fmt.Println(input.PaymentType)

	err = h.paymentService.ProcessPayment(input)
	if err != nil {
		response := helper.ResponseJSON("Failed to process notification", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, input)
}