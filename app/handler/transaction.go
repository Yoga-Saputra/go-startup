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
		response := helper.ApiResponse(err.Error(), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("Campaign detail transaction", http.StatusOK, "success", transaction.FormatCampaignSlice(transactions))
	c.JSON(http.StatusOK, response)
}
