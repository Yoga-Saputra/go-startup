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
