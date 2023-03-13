package transaction

import "time"

type CampaignTransactionFormatter struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// slice transactions
func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {

	formater := CampaignTransactionFormatter{
		Id:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}

	return formater
}

// mapping  slice campaigns
func FormatCampaignSlice(transactions []Transaction) []CampaignTransactionFormatter {
	campaignsFormatter := []CampaignTransactionFormatter{}

	if len(transactions) == 0 {
		return campaignsFormatter
	}

	for _, transaction := range transactions {
		transactionFormatter := FormatCampaignTransaction(transaction)
		campaignsFormatter = append(campaignsFormatter, transactionFormatter)
	}

	return campaignsFormatter
}

type UserTransactionFormater struct {
	Id        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}

type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

func FormatUserTransaction(transaction Transaction) UserTransactionFormater {
	formatter := UserTransactionFormater{}
	formatter.Id = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	campaignFormatter := CampaignFormatter{}
	campaignFormatter.Name = transaction.Campaign.Name
	campaignFormatter.ImageURL = ""

	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	formatter.Campaign = campaignFormatter

	return formatter
}

// mapping  slice users
func FormatUserSlice(transactions []Transaction) []UserTransactionFormater {
	usersFormatter := []UserTransactionFormater{}

	if len(transactions) == 0 {
		return usersFormatter
	}

	for _, transaction := range transactions {
		transactionFormatter := FormatUserTransaction(transaction)
		usersFormatter = append(usersFormatter, transactionFormatter)
	}

	return usersFormatter
}

type TransactionFormatter struct {
	Id         int    `json:"id"`
	CampaignID int    `json:"name"`
	UserID     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

// transactions
func FormatTransaction(transaction Transaction) TransactionFormatter {

	formater := TransactionFormatter{
		Id:         transaction.ID,
		CampaignID: transaction.CampaignID,
		Amount:     transaction.Amount,
		UserID:     transaction.UserID,
		Status:     transaction.Status,
		Code:       transaction.Code,
		PaymentURL: transaction.PaymentURL,
	}

	return formater
}
