package transaction

import (
	"errors"
	"startup/app/campaign"
	"startup/config"
	"strconv"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionByCampaignId(input GetCampaignTransactionInput) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionByCampaignId(input GetCampaignTransactionInput) ([]Transaction, error) {

	campaign, err := s.campaignRepository.FindById(input.ID)
	config.Loggers("error", campaign)
	if err != nil {
		return []Transaction{}, err
	}

	msg := "no data of campaign with campaign id = " + strconv.Itoa(input.ID)
	if campaign.UserID == 0 {
		return []Transaction{}, errors.New(msg)
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("not an owner of the campaign")
	}

	transaction, err := s.repository.GetByCampaignId(input.ID)
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
