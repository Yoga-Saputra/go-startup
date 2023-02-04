package handler

import (
	"net/http"
	"startup/app/campaign"
	"startup/app/helper"
	"startup/app/users"
	"startup/config"
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
	ErrorResponseCampaign(ctx, err, "failed to get detail of campaign")

	campaign, err := ch.service.GetCampaignById(input)

	ErrorResponseCampaign(ctx, err, "failed to get detail of campaign")
	SuccessResponseCampaign(ctx, campaign)
}

func (ch *campaignHandler) CreateCampaign(ctx *gin.Context) {
	var input campaign.CreateCampaignInput

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, ctx, "create campaign failed", "error", http.StatusBadRequest, errors)
		return
	}

	currentUser := ctx.MustGet("current_user").(users.User)
	input.User = currentUser
	newCampaign, err := ch.service.CreateCampaign(input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		helper.ErrorValidation(err, ctx, "create campaign failed", "error", http.StatusBadRequest, errors)
		return
	}

	response := helper.ApiResponse("campaign detail", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	ctx.JSON(http.StatusOK, response)
}

func (ch *campaignHandler) UpdateCampaign(ctx *gin.Context) {
	var inputId campaign.GetCampaignDetailInput

	err := ctx.ShouldBindUri(&inputId)

	ErrorResponseCampaign(ctx, err, "failed to update campaign 123")

	var inputData campaign.CreateCampaignInput

	err = ctx.ShouldBindJSON(&inputData)

	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"error": errors}
		helper.ErrorValidation(err, ctx, "failed to update campaign 3", "error", http.StatusBadRequest, errorMessage)
		return
	}

	currentUser := ctx.MustGet("current_user").(users.User)
	inputData.User = currentUser

	updateCampaign, err := ch.service.UpdateCampaign(inputId, inputData)
	if err != nil {
		ErrorResponseCampaign(ctx, err, "failed to update campaign")
		return
	}

	response := helper.ApiResponse("success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updateCampaign))
	ctx.JSON(http.StatusOK, response)

}

func ErrorResponseCampaign(ctx *gin.Context, err error, msg string) {
	if err != nil {
		config.Loggers("error", err)
		response := helper.ApiResponse(msg, http.StatusBadRequest, "error", err.Error())
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
