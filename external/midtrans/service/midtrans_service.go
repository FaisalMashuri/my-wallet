package service

import (
	"github.com/FaisalMashuri/my-wallet/config"
	midtrans_ext "github.com/FaisalMashuri/my-wallet/external/midtrans"
	"github.com/FaisalMashuri/my-wallet/internal/domain/topup"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type midtransService struct {
	client snap.Client
}

func (m *midtransService) GenerateSnapURL(t *topup.TopUp) error {
	//TODO implement me
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  t.ID,
			GrossAmt: int64(t.Amount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	snapResp, err := m.client.CreateTransaction(req)
	if err != nil {
		return err
	}
	t.SnapURL = snapResp.RedirectURL
	return nil
}

func (m *midtransService) VerifyPayment(orderId string) (bool, error) {
	var client coreapi.Client
	envMidtrans := midtrans.Sandbox
	if config.AppConfig.Midtrans.Env == "production" {
		envMidtrans = midtrans.Sandbox
	}
	client.New(config.AppConfig.Midtrans.ServerKey, envMidtrans)

	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := client.CheckTransaction(orderId)
	if e != nil {
		return false, e
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// TODO set transaction status on your databaase to 'success'
				return true, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}
	return true, nil
}

func NewService() midtrans_ext.MidtransService {
	var snap snap.Client
	envMidtrans := midtrans.Sandbox
	if config.AppConfig.Midtrans.Env == "production" {
		envMidtrans = midtrans.Production
	}
	midtrans.ServerKey = config.AppConfig.Midtrans.ServerKey
	snap.New(config.AppConfig.Midtrans.ServerKey, envMidtrans)
	return &midtransService{
		snap,
	}
}
