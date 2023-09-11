package service

import (
	"fmt"
	midtrans_ext "github.com/FaisalMashuri/my-wallet/external/midtrans"
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	notifResponse "github.com/FaisalMashuri/my-wallet/internal/domain/notification/dto/response"
	"github.com/FaisalMashuri/my-wallet/internal/domain/sse/dto"
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup"
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup/dto/response"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"log"
)

type topUpService struct {
	hub                 *dto.Hub
	topUpRepository     topup.TopUpRepository
	midtransService     midtrans_ext.MidtransService
	notificationService notification.NotificationRepository
	accountRepository   account.AccountRepository
}

func (t *topUpService) NotificationAfterTransfer(account account.Account, amount float64) {
	//TODO implement me
	fmt.Println("AMONT TOPUP :  ", amount)
	notificationSender := notification.Notification{
		UserID: account.UserID,
		Title:  "Transfer Berhasil",
		Body:   fmt.Sprintf("Topup senilai %.2f berhasil", amount),
		IsRead: 0,
	}
	err := t.notificationService.InsertNotification(&notificationSender)
	if err != nil {
		log.Println("Error notif sender : ", err.Error())
	}
	fmt.Println("Insert notif berhasil")
	fmt.Println("Pertama : ", t.hub.NotificationChannel[account.UserID])
	fmt.Println("KEDUA : ", t.hub.NotificationChannel)
	channel, ok := t.hub.NotificationChannel[account.UserID]
	fmt.Println("CHANNEL : ", channel)
	fmt.Println("OK BANGET  : ", ok)
	if ok {
		channel <- notifResponse.NotificationDataRes{
			ID:        notificationSender.ID,
			UserID:    account.UserID,
			Title:     notificationSender.Title,
			Body:      notificationSender.Body,
			IsRead:    notificationSender.IsRead,
			CreatedAt: notificationSender.CreatedAt,
			Status:    notificationSender.Status,
		}
	}

}

func (t topUpService) ConfirmedTopUp(id string) error {
	//TODO implement me
	topUp, err := t.topUpRepository.FindById(id)
	if err != nil {
		return err
	}
	if topUp == nil {
		return errors.New(contract.ErrRecordNotFound)
	}
	fmt.Println()
	accountData, err := t.accountRepository.FindAccountByUserId(topUp.UserID)
	if err != nil {
		return err
	}
	if accountData == nil {
		return errors.New(contract.ErrRecordNotFound)
	}
	accountData.Balance += topUp.Amount
	_, err = t.accountRepository.UpdateBalance(*accountData)
	if err != nil {
		return err
	}

	fmt.Println(accountData)
	go t.NotificationAfterTransfer(*accountData, topUp.Amount)
	return nil
}

func (t topUpService) InitializeTopUp(req request.TopUpRequest) (response.TopUpResponnse, error) {
	//TODO implement me
	topUp := topup.TopUp{
		ID:     uuid.New().String(),
		UserID: req.UserID,
		Status: 0,
		Amount: req.Amount,
	}
	err := t.midtransService.GenerateSnapURL(&topUp)
	if err != nil {
		return response.TopUpResponnse{}, err
	}
	err = t.topUpRepository.Insert(&topUp)
	if err != nil {
		return response.TopUpResponnse{}, err
	}
	return response.TopUpResponnse{
		SnapURL: topUp.SnapURL,
	}, nil
}

func NewService(topUpRepository topup.TopUpRepository, midtransService midtrans_ext.MidtransService, notificationService notification.NotificationRepository, accountRepository account.AccountRepository, hub *dto.Hub) topup.TopUpService {
	return &topUpService{
		topUpRepository:     topUpRepository,
		midtransService:     midtransService,
		notificationService: notificationService,
		accountRepository:   accountRepository,
		hub:                 hub,
	}
}
