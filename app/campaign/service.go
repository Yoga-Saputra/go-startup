package campaign

type Service interface {
	GetCampains(userID int) ([]Campaign, error)
	GetCampaignById(input GetCampaignDetailInput) (Campaign, error)
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

func returnHelper(err error, campaign []Campaign) ([]Campaign, error) {
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
