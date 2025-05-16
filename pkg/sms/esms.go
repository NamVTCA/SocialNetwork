package sms

// import (
// 	"fmt"

// 	"socialnetwork/pkg/utils"
// 	"github.com/go-resty/resty/v2"
// )

// type ESMSender struct {
// 	ApiKey    string
// 	SecretKey string
// 	BrandName string
// }

// func NewESMSender(apiKey, secretKey, brandName string) SMSSender {
// 	return &ESMSender{
// 		ApiKey:    apiKey,
// 		SecretKey: secretKey,
// 		BrandName: "",
// 	}
// }

// func (e *ESMSender) Send(to, body string) error {
// 	client := resty.New()
// 	to = utils.FormatPhoneToVietnamese(to)

// 	resp, err := client.R().
// 		SetQueryParams(map[string]string{
// 			"ApiKey":    e.ApiKey,
// 			"SecretKey": e.SecretKey,
// 			"Phone":     to,
// 			"Content":   body,
// 			"Brandname": "",    // để trống nếu chưa đăng ký brandname
// 			"SmsType":   "1",   // gửi SMS thường
// 		}).
// 		Get("https://api.esms.vn/MainService.svc/json/SendMultipleMessage_V4_get")

// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("📤 eSMS response:", resp.String())
// 	return nil
// }



