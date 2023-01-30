package handler

import (
	"net/http"
	"startup/app/campaign"
	"startup/app/helper"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (ch *campaignHandler) FindAllCamp(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("user_id"))
	camp, err := ch.service.GetCampains(userId)

	if err != nil {
		response := helper.ApiResponse("error to get campaigns", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
	}
	response := helper.ApiResponse("list of campaigns", http.StatusOK, "success", camp)

	ctx.JSON(http.StatusOK, response)
}
