package sms

// import (
// 	"fmt"
// 	"github.com/twilio/twilio-go"
// 	openapi "github.com/twilio/twilio-go/rest/api/v2010"
// )

// type TwilioSender struct {
// 	client  *twilio.RestClient
// 	fromNum string
// }

// func NewTwilioSender(accountSID, authToken, fromNum string) SMSSender {
// 	client := twilio.NewRestClientWithParams(twilio.ClientParams{
// 		Username: accountSID,
// 		Password: authToken,
// 	})
// 	return &TwilioSender{
// 		client:  client,
// 		fromNum: fromNum,
// 	}
// }

// func (t *TwilioSender) Send(to string, body string) error {
//     fmt.Printf("[DEBUG] TwilioSender.Send called\n")
//     fmt.Printf("[DEBUG] To: %s\n", to)
//     fmt.Printf("[DEBUG] From: %s\n", t.fromNum)
//     fmt.Printf("[DEBUG] Body: %s\n", body)

//     params := &openapi.CreateMessageParams{}
//     params.SetTo(to)
//     params.SetFrom(t.fromNum)
//     params.SetBody(body)

//     resp, err := t.client.Api.CreateMessage(params)
//     if err != nil {
//         fmt.Printf("[ERROR] Twilio API error: %v\n", err)
//         return err
//     }

//     if resp.Sid != nil {
//         fmt.Printf("[INFO] Twilio SMS sent successfully. SID: %s\n", *resp.Sid)
//     } else {
//         fmt.Printf("[WARN] Twilio SMS sent but no SID received\n")
//     }
//     return nil
// }
