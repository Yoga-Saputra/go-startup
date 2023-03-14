package transaction

import (
	"errors"
	"fmt"
	"startup/app/campaign"
	"startup/app/payment"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/google/uuid"
)

type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService     payment.Service
}

type Service interface {
	GetTransactionByCampaignId(input GetCampaignTransactionInput) ([]Transaction, error)
	GetTransactionByUserId(userId int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
	ExportExcel(data []UserTransactionExcel) error
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionByCampaignId(input GetCampaignTransactionInput) ([]Transaction, error) {

	campaign, err := s.campaignRepository.FindById(input.ID)
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

func (s *service) GetTransactionByUserId(userId int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserId(userId)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	uuID := uuid.New()
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "pending"
	transaction.Code = uuID.String()

	newTransaction, err := s.repository.Save(transaction)

	if err != nil {
		return newTransaction, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTransaction.ID,
		Amount: newTransaction.Amount,
		Code:   newTransaction.Code,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)

	if err != nil {
		return newTransaction, err
	}

	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(newTransaction)

	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {

	transaction, err := s.repository.GetById(input.OrderID)

	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	}

	if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	}

	if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)

	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindById(updatedTransaction.CampaignID)

	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err = s.campaignRepository.Update(campaign)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) ExportExcel(transaction []UserTransactionExcel) error {
	headers := map[string]string{
		"A1": "transaction_id",
		"B1": "amount",
		"C1": "status",
		"D1": "transaction_code",
		"E1": "campaign_name",
		"F1": "campaign_image_url",
		"G1": "created_at",
	}
	file := excelize.NewFile()
	for k, v := range headers {
		file.SetCellValue("Sheet1", k, v)
	}
	fmt.Println(transaction)
	for i := 0; i < len(transaction); i++ {
		appendRow(file, i, transaction)
	}

	var filename string = fmt.Sprintf("storage/excel/transaction.xlsx")

	err := file.SaveAs(filename)
	if err != nil {
		return err
	}

	return nil
}

// Append every row, we can add styles if neeeded
func appendRow(file *excelize.File, index int, transaction []UserTransactionExcel) (fileWriter *excelize.File) {
	rowCount := index + 2
	file.SetCellValue("Sheet1", fmt.Sprintf("A%v", rowCount), transaction[index].Id)
	file.SetCellValue("Sheet1", fmt.Sprintf("B%v", rowCount), transaction[index].Amount)
	file.SetCellValue("Sheet1", fmt.Sprintf("C%v", rowCount), transaction[index].Status)
	file.SetCellValue("Sheet1", fmt.Sprintf("D%v", rowCount), transaction[index].TransactionId)
	file.SetCellValue("Sheet1", fmt.Sprintf("E%v", rowCount), transaction[index].CampaignName)
	file.SetCellValue("Sheet1", fmt.Sprintf("F%v", rowCount), transaction[index].CampaignImageUrl)
	file.SetCellValue("Sheet1", fmt.Sprintf("G%v", rowCount), transaction[index].CreatedAt)

	return file
}
