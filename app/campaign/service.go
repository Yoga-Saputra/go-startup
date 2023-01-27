package campaign

type Service interface {
	GetCampains(userID int) ([]Campaign, error)
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
		return errorHelper(err, campaign)
	}

	campaign, err := s.repository.FindAll()
	return errorHelper(err, campaign)
}

func errorHelper(err error, campaign []Campaign) ([]Campaign, error) {
	if err != nil {
		return campaign, err
	}

	return campaign, nil
}
