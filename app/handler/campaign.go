package handler

import (
	"fmt"
	"net/http"
	"startup/app/campaign"
	"startup/app/helper"
	"startup/app/users"
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

	if err != nil {
		ErrorResponseCampaign("failed to get detail of campaign", ctx, err, err.Error())
		return
	}

	campaign, err := ch.service.GetCampaignById(input)

	if err != nil {
		ErrorResponseCampaign("failed to get detail of campaign", ctx, err, err.Error())
		return
	}

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

	if err != nil {
		ErrorResponseCampaign("failed to update campaign", ctx, err, err.Error())
		return
	}

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
		ErrorResponseCampaign("failed to update campaign", ctx, err, err.Error())
		return
	}

	response := helper.ApiResponse("success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updateCampaign))
	ctx.JSON(http.StatusOK, response)

}

func (ch *campaignHandler) UploadImage(ctx *gin.Context) {
	var input campaign.CreateCampaignImageInput
	err := ctx.ShouldBind(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)
		response := helper.ApiResponse("failed to upload campaign image", http.StatusBadRequest, "error", errors)

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	checkCampaign, err := ch.service.CheckCampaignService(input)

	if checkCampaign.ID == 0 {
		stringId := strconv.Itoa(input.CampaignId)
		msg := "no data of campaign with campaign id " + stringId
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse(msg, http.StatusBadRequest, "error", data)

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if err != nil {
		data := gin.H{"is_uploaded": false}
		ErrorResponseCampaign("failed to upload campaign image", ctx, err, data)
		return
	}

	file, err := ctx.FormFile("file")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		ErrorResponseCampaign("failed to upload campaign image", ctx, err, data)
		return
	}

	currentUser := ctx.MustGet("current_user").(users.User)
	userId := currentUser.ID
	input.User.ID = userId
	path := fmt.Sprintf("storage/images/campaigns/%d-%s", userId, file.Filename)

	// save image to derectory
	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		ErrorResponseCampaign("failed to upload campaign image", ctx, err, data)
		return
	}

	// save image to database
	_, err = ch.service.SaveCampaignImage(input, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		ErrorResponseCampaign("failed to upload campaign image", ctx, err, data)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("avatar successfully uploaded", http.StatusOK, "success", data)

	ctx.JSON(http.StatusOK, response)
}

func ErrorResponseCampaign(msg string, ctx *gin.Context, err error, data interface{}) {
	response := helper.ApiResponse(msg, http.StatusBadRequest, "error", data)
	ctx.JSON(http.StatusBadRequest, response)
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
