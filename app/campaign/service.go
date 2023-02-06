package campaign

import (
	"errors"
	"fmt"
	"startup/config"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampains(userID int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputId GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
	CheckCampaignService(input CreateCampaignImageInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampains(userId int) ([]Campaign, error) {
	if userId != 0 {
		campaign, err := s.repository.FindByUserId(userId)
		return returnHelper(err, campaign)
	}

	campaign, err := s.repository.FindAll()
	return returnHelper(err, campaign)
}

func (s *service) GetCampaignById(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindById(input.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CheckCampaignService(input CreateCampaignImageInput) (Campaign, error) {
	campaign, err := s.repository.CheckCampaignRepository(input.CampaignId)
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	// create slug
	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)

	campaign := Campaign{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		UserID:           input.User.ID,
		Slug:             slug.Make(slugCandidate),
	}

	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

// update campaign
func (s *service) UpdateCampaign(inputId GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindById(inputId.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updateCampaign, err := s.repository.Update(campaign)

	if err != nil {
		return updateCampaign, err
	}

	return updateCampaign, nil

}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1

		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignId)

		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaign, err := s.repository.FindById(input.CampaignId)

	if err != nil {
		return CampaignImage{}, err
	}

	mapLog := map[string]interface{}{
		"campaign_user_id": campaign.UserID,
		"userId":           input.User.ID,
	}
	config.Loggers("error", mapLog)

	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("you dont have permission for this campaign")
	}

	campaignImage := CampaignImage{
		CampaignID: input.CampaignId,
		IsPrimary:  isPrimary,
		FileName:   fileLocation,
	}

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	config.Loggers("error", err)

	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, nil
}

func returnHelper(err error, campaign []Campaign) ([]Campaign, error) {
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
