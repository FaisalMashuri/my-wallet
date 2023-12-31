package service

import (
	"encoding/json"
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	"github.com/FaisalMashuri/my-wallet/internal/domain/notification"
	notifResponse "github.com/FaisalMashuri/my-wallet/internal/domain/notification/dto/response"
	"github.com/FaisalMashuri/my-wallet/internal/domain/sse/dto"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/response"
	"github.com/FaisalMashuri/my-wallet/internal/domain/user"
	"github.com/redis/go-redis/v9"

	"github.com/FaisalMashuri/my-wallet/shared"
	"github.com/FaisalMashuri/my-wallet/shared/contract"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"log"
	"time"
)

type transactionService struct {
	repoTransaction transaction.TransactionRepository
	repoAccount     account.AccountRepository
	repoNotif       notification.NotificationRepository
	hub             *dto.Hub
	redisClient     *redis.Client
}

func (t *transactionService) NotificationAfterTransfer(sofAccount account.Account, dofAccount account.Account, amount float64) {
	//TODO implement me
	notificationSender := notification.Notification{
		UserID: sofAccount.UserID,
		Title:  "Transfer Berhasil",
		Body:   fmt.Sprintf("Transfer senilai %.2f kepada %s berhasil", amount, dofAccount.AccountNumber),
		IsRead: 0,
	}
	err := t.repoNotif.InsertNotification(&notificationSender)
	if err != nil {
		log.Println("Error notif sender : ", err.Error())
	}
	if channel, ok := t.hub.NotificationChannel[sofAccount.UserID]; ok {
		channel <- notifResponse.NotificationDataRes{
			ID:        notificationSender.ID,
			UserID:    sofAccount.UserID,
			Title:     notificationSender.Title,
			Body:      notificationSender.Body,
			IsRead:    notificationSender.IsRead,
			CreatedAt: notificationSender.CreatedAt,
			Status:    notificationSender.Status,
		}
	}
	notificationReciever := notification.Notification{
		UserID: dofAccount.UserID,
		Title:  "Transfer Diterima",
		Body:   fmt.Sprintf("Transfer senilai %.2f diterima dari %s", amount, sofAccount.AccountNumber),
		IsRead: 0,
	}
	err = t.repoNotif.InsertNotification(&notificationReciever)
	if err != nil {
		log.Println("Error notif reciever : ", err.Error())
	}
	if channel, ok := t.hub.NotificationChannel[dofAccount.UserID]; ok {
		channel <- notifResponse.NotificationDataRes{
			ID:        notificationReciever.ID,
			UserID:    dofAccount.UserID,
			Title:     notificationReciever.Title,
			Body:      notificationReciever.Body,
			IsRead:    notificationReciever.IsRead,
			CreatedAt: notificationReciever.CreatedAt,
			Status:    notificationReciever.Status,
		}
	}
}

func (t transactionService) TranferInquiry(InquiryReq request.TransferInquiryReq, ctx *fiber.Ctx) (*response.TransferInquiryRes, error) {
	//TODO implement me

	credentialuser := ctx.Locals("credentials").(user.User)
	myAccount, err := t.repoAccount.FindAccountByAccountNumber(InquiryReq.SofAccountNumber)
	if myAccount == nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(contract.ErrRecordNotFound)
	}
	if myAccount.UserID != credentialuser.ID {
		return nil, errors.New(contract.ErrTransactionUnauthorized)
	}

	dofAccount, err := t.repoAccount.FindAccountByAccountNumber(InquiryReq.DofAccountNumber)
	if dofAccount == nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(contract.ErrRecordNotFound)
	}

	if myAccount.Balance < InquiryReq.Amount {
		return nil, errors.New(contract.ErrInsufficientBalance)
	}
	inquiryKey := shared.GenerateInquiryKey()
	InquiryKeyRes := request.TransferInquiryReq{
		Amount:           InquiryReq.Amount,
		SofAccountNumber: InquiryReq.SofAccountNumber,
		DofAccountNumber: InquiryReq.DofAccountNumber,
	}

	inquiryJSONString, err := json.Marshal(InquiryKeyRes)

	inquirymodel := transaction.TransactionInquiry{
		InquiryKey: inquiryKey,
		Value:      string(inquiryJSONString),
	}
	//
	//data, err := t.repoTransaction.CreateTransactionInquiry(inquirymodel)
	//if err != nil {
	//	return nil, err
	//}

	inquiryRespnse := response.TransferInquiryRes{
		inquirymodel.InquiryKey,
	}
	err = t.redisClient.Set(ctx.Context(), inquirymodel.InquiryKey, inquirymodel.Value, 30*time.Minute).Err()
	if err != nil {
		log.Println(err)
	}

	return &inquiryRespnse, nil
}

func (t transactionService) TransferInquiryExec(InquiryExecReq request.TransferInquiryExec, ctx *fiber.Ctx) error {
	//TODO implement me
	var inqReq request.TransferInquiryReq

	//Get data from redis
	val, err := t.redisClient.Get(ctx.Context(), InquiryExecReq.InquiryKey).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println(err.Error())
		}
		return errors.New(contract.ErrRecordNotFound)

	}
	fmt.Println("data val redis  : ", val)
	err = json.Unmarshal([]byte(val), &inqReq)
	if err != nil {
		log.Println("Error : ", err.Error())
	}
	fmt.Println("data dari redis : ", inqReq)

	defer func() {
		_ = t.redisClient.Del(ctx.Context(), InquiryExecReq.InquiryKey).Err()
	}()

	myAccount, err := t.repoAccount.FindAccountByAccountNumber(inqReq.SofAccountNumber)
	if myAccount == nil {
		if err != nil {
			log.Println("Error Find account by user id : ", err.Error())

			return err
		}
		return errors.New(contract.ErrRecordNotFound)
	}

	dofAccount, err := t.repoAccount.FindAccountByAccountNumber(inqReq.DofAccountNumber)
	if err != nil {
		log.Println("Error Find Account: ", err.Error())

		return err
	}

	debitTransaction := transaction.Transaction{
		AccountID:       myAccount.ID,
		SofNumber:       myAccount.AccountNumber,
		DofNumber:       dofAccount.AccountNumber,
		TransactionType: "D",
		Amount:          inqReq.Amount,
		TransactionTime: time.Now(),
	}
	_, err = t.repoTransaction.CreateTransaction(debitTransaction)
	if err != nil {
		log.Println("Error debit transaction : ", err.Error())
		return err
	}

	creditTransaction := transaction.Transaction{
		AccountID:       dofAccount.ID,
		SofNumber:       dofAccount.AccountNumber,
		DofNumber:       myAccount.AccountNumber,
		TransactionType: "C",
		Amount:          inqReq.Amount,
		TransactionTime: time.Now(),
	}
	_, err = t.repoTransaction.CreateTransaction(creditTransaction)
	if err != nil {
		log.Println("Error credit : ", err.Error())

		return err
	}

	fmt.Println("my account : ", myAccount)
	fmt.Println("My Account Balance : ", myAccount.Balance)
	fmt.Println("DOF ACCount Balance : ", dofAccount.Balance)

	fmt.Println("DOF Account : ", dofAccount)

	myAccount.Balance = myAccount.Balance - inqReq.Amount
	_, err = t.repoAccount.UpdateBalance(*myAccount)
	if err != nil {
		return err
	}

	dofAccount.Balance += inqReq.Amount
	_, err = t.repoAccount.UpdateBalance(*dofAccount)
	if err != nil {
		return err
	}
	err = t.repoTransaction.DeleteInquiry(InquiryExecReq.InquiryKey)
	if err != nil {
		log.Println("Error delete : ", err.Error())

		return err
	}
	go t.NotificationAfterTransfer(*myAccount, *dofAccount, inqReq.Amount)
	return nil
}

func NewService(repoTransaction transaction.TransactionRepository, repoAccount account.AccountRepository, repoNotif notification.NotificationRepository, hub *dto.Hub, redisClient *redis.Client) transaction.TransactionService {
	return &transactionService{
		repoTransaction: repoTransaction,
		repoAccount:     repoAccount,
		repoNotif:       repoNotif,
		hub:             hub,
		redisClient:     redisClient,
	}
}
