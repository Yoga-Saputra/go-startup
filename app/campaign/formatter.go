package campaign

import "strings"

// helper campaign
type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amount"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

type CampaigUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type CampaigImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

type CampaignDetailFormatter struct {
	ID               int                     `json:"id"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"short_description"`
	Description      string                  `json:"description"`
	ImageUrl         string                  `json:"image_url"`
	GoalAmount       int                     `json:"goal_amount"`
	CurrentAmount    int                     `json:"current_amount"`
	BackerCount      int                     `json:"backers_count"`
	UserID           int                     `json:"user_id"`
	Slug             string                  `json:"slug"`
	Perks            []string                `json:"perks"`
	User             CampaigUserFormatter    `json:"user"`
	Image            []CampaigImageFormatter `json:"images"`
}

// single campaign
func FormatCampaign(campaign Campaign) CampaignFormatter {
	var image_url string

	if len(campaign.CampaignImages) > 0 {
		image_url = campaign.CampaignImages[0].FileName
	}

	formater := CampaignFormatter{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		ImageURL:         image_url,
		Slug:             campaign.Slug,
	}

	return formater
}

// mapping  slice campaigns
func FormatCampaignSlice(campaigns []Campaign) []CampaignFormatter {
	campaignsFormatter := []CampaignFormatter{}

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaign(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}

func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	var image_url string
	var perks []string

	if len(campaign.CampaignImages) > 0 {
		image_url = campaign.CampaignImages[0].FileName
	}
	// break the string  with strings.Split(campaign.Perks, ","),
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk)) //delete space string and combine them
	}
	user := campaign.User
	campaignUserFormatter := CampaigUserFormatter{
		Name:     user.Name,
		ImageURL: user.AvatarFileName,
	}

	images := []CampaigImageFormatter{}
	for _, image := range campaign.CampaignImages {
		isPrimary := false
		if image.IsPrimary == 1 {
			isPrimary = true
		}

		campaignImageFormatter := CampaigImageFormatter{
			ImageURL:  image.FileName,
			IsPrimary: isPrimary,
		}

		images = append(images, campaignImageFormatter)
	}

	campaignDetailFormatter := CampaignDetailFormatter{
		ID:               campaign.ID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		ImageUrl:         image_url,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		BackerCount:      campaign.BackerCount,
		UserID:           campaign.UserID,
		Slug:             campaign.Slug,
		Perks:            perks,
		User:             campaignUserFormatter,
		Image:            images,
	}

	return campaignDetailFormatter
}

// mapping slice campaign detail
func FormatCampaignDetailSlice(campaigns []Campaign) []CampaignDetailFormatter {
	campaignsFormatter := []CampaignDetailFormatter{}

	for _, campaign := range campaigns {
		campaignFormatter := FormatCampaignDetail(campaign)
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}

	return campaignsFormatter
}
