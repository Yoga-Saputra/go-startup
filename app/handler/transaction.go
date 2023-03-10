package handler

import (
	"net/http"
	"startup/app/helper"
	"startup/app/transaction"
	"startup/app/users"

	"github.com/gin-gonic/gin"
)

type transactionnhandler struct {
	sercive transaction.Service
}

func NewTransaction(service transaction.Service) *transactionnhandler {
	return &transactionnhandler{service}
}

func (h *transactionnhandler) GetCampaignTransaction(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.ApiResponse(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("current_user").(users.User)

	input.User = currentUser

	transactions, err := h.sercive.GetTransactionByCampaignId(input)
	if err != nil {
		response := helper.ApiResponse("failed to get campaign's transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("Campaign's transaction", http.StatusOK, "success", transaction.FormatCampaignSlice(transactions))
	c.JSON(http.StatusOK, response)
}

func (h *transactionnhandler) GetUserTransaction(c *gin.Context) {
	currentUser := c.MustGet("current_user").(users.User)

	userID := currentUser.ID

	transactions, err := h.sercive.GetTransactionByUserId(userID)

	if err != nil {
		response := helper.ApiResponse("failed to get user's transaction", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("User's transaction", http.StatusOK, "success", transaction.FormatUserSlice(transactions))
	c.JSON(http.StatusOK, response)

}

func (h *transactionnhandler) ExcelTransaction(c *gin.Context) {
	currentUser := c.MustGet("current_user").(users.User)

	userID := currentUser.ID

	data, err := h.sercive.GetTransactionByUserId(userID)

	if err != nil {
		response := helper.ApiResponse(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	trans := transaction.FormatUserSliceExcel(data)
	// generate excel
	err = h.sercive.ExportExcel(trans)
	if err != nil {
		response := helper.ApiResponse(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if len(data) > 0 {
		response := helper.ApiResponse("successfully export transaction to excel", http.StatusOK, "success", trans)
		c.JSON(http.StatusOK, response)
		return
	}

	response := helper.ApiResponse("no data found", http.StatusNotFound, "false", nil)
	c.JSON(http.StatusOK, response)

}

func (h *transactionnhandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, c, "failed to create transaction", "error", http.StatusBadRequest, errors)
		return
	}

	currentUser := c.MustGet("current_user").(users.User)
	input.User = currentUser

	newTransaction, err := h.sercive.CreateTransaction(input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, c, "failed to create transaction", "error", http.StatusBadRequest, errors)
		return
	}

	response := helper.ApiResponse("success create transaction", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))
	c.JSON(http.StatusOK, response)

}

func (h *transactionnhandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		helper.ErrorValidation(err, c, "failed to process transaction", "error", http.StatusBadRequest, nil)
		return
	}

	err = h.sercive.ProcessPayment(input)
	if err != nil {
		helper.ErrorValidation(err, c, "failed to process transaction", "error", http.StatusBadRequest, nil)
		return
	}

	c.JSON(http.StatusOK, input)
}
