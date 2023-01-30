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

func (ch *campaignHandler) GetAllCamp(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("user_id"))
	camp, err := ch.service.GetCampains(userId)

	if err != nil {
		response := helper.ApiResponse("error to get campaigns", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
	}
	response := helper.ApiResponse("list of campaigns", http.StatusOK, "success", campaign.FormatCampaignSlice(camp))

	ctx.JSON(http.StatusOK, response)
}

func (ch *campaignHandler) GetCampain(ctx *gin.Context) {
	var input campaign.GetCampaignDetailInput

	err := ctx.ShouldBindUri(&input)
	ErrorResponseCampaign(ctx, err)

	campaign, err := ch.service.GetCampaignById(input)

	ErrorResponseCampaign(ctx, err)
	SuccessResponseCampaign(ctx, campaign)
}

func ErrorResponseCampaign(ctx *gin.Context, err error) {
	if err != nil {
		response := helper.ApiResponse("failed to get detail of campaign", http.StatusBadRequest, "error", err)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
}

func SuccessResponseCampaign(ctx *gin.Context, campaignDetail campaign.Campaign) {
	if campaignDetail.ID == 0 {
		arrNull := []campaign.Campaign{}
		response := helper.ApiResponse("no detail data of campaign", http.StatusBadRequest, "error", campaign.FormatCampaignDetailSlice(arrNull))
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.ApiResponse("campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	ctx.JSON(http.StatusOK, response)
}
