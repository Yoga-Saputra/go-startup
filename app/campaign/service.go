package campaign

import (
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampains(userID int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
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

func returnHelper(err error, campaign []Campaign) ([]Campaign, error) {
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
