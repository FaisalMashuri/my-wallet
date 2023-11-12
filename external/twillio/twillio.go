package twillio

import (
	"fmt"
	"github.com/twilio/twilio-go"
	openApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendWhatsapp(phoneNumber string, otp string) error {
	client := twilio.NewRestClient()
	params := &openApi.CreateMessageParams{}
	params.SetTo(fmt.Sprintf("whatsapp:+62%s", phoneNumber))
	params.SetFrom("whatsapp:+14155238886")

	params.SetBody(fmt.Sprintf("Jangan berikan kode otp ini ke siapapun termasuk dari pihak kami.\n Berikut kode OTP %s", otp))

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return err
	}
	fmt.Println("WA SENT SUCCESFULL")
	return nil

}
