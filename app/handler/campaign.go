package handler

import (
	"net/http"
	"startup/app/campaign"
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
		ctx.JSON(http.StatusBadRequest, err.Error())
	}

	ctx.JSON(http.StatusOK, camp)
}
