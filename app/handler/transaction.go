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
