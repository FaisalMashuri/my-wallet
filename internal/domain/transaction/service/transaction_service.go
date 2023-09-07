package service

import (
	"encoding/json"
	"fmt"
	"github.com/FaisalMashuri/my-wallet/internal/domain/account"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/request"
	"github.com/FaisalMashuri/my-wallet/internal/domain/transaction/dto/response"
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
}

func (t transactionService) TranferInquiry(InquiryReq request.TransferInquiryReq, ctx *fiber.Ctx) (*response.TransferInquiryRes, error) {
	//TODO implement me
	//credentialuser := ctx.Locals("credentials").(user.User)
	myAccount, err := t.repoAccount.FindAccountByAccountNumber(InquiryReq.SofAccountNumber)
	if myAccount == nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(contract.ErrRecordNotFound)
	}

	dofAccount, err := t.repoAccount.FindAccountByAccountNumber(InquiryReq.DofAccountNumber)
	if dofAccount == nil {
		if err != nil {
			return nil, err
		}
		return nil, errors.New(contract.ErrRecordNotFound)
	}

	if myAccount.Balance < InquiryReq.Amount {
		return nil, errors.New(contract.ErrInsuficentBalance)
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

	data, err := t.repoTransaction.CreateTransactionInquiry(inquirymodel)
	if err != nil {
		return nil, err
	}

	inquiryRespnse := response.TransferInquiryRes{
		data.InquiryKey,
	}

	return &inquiryRespnse, nil
}

func (t transactionService) TransferInquiryExec(InquiryExecReq request.TransferInquiryExec, ctx *fiber.Ctx) error {
	//TODO implement me

	dataInquiry, err := t.repoTransaction.FindTransactionInquiry(InquiryExecReq.InquiryKey)
	if err != nil {
		log.Println("Error  Find Inquiry: ", err.Error())
		return err
	}
	fmt.Println("data inquiry : ", dataInquiry.Value)

	var inqReq request.TransferInquiryReq
	err = json.Unmarshal([]byte(dataInquiry.Value), &inqReq)
	if err != nil {
		log.Println("Error Unmarshar Inquiry Val : ", err.Error())
		return err
	}
	fmt.Println("inreq: ", inqReq)

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
	err = t.repoTransaction.DeleteInquiry(dataInquiry.InquiryKey)
	if err != nil {
		log.Println("Error delete : ", err.Error())

		return err
	}
	return nil
}

func NewService(repoTransaction transaction.TransactionRepository, repoAccount account.AccountRepository) transaction.TransactionService {
	return &transactionService{
		repoTransaction: repoTransaction,
		repoAccount:     repoAccount,
	}
}
