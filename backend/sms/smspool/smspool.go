package smspool

import (
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/bogdanfinn/tls-client/profiles"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/tidwall/gjson"
)

func NewSMSPoolClient(apiKey string) (SMSPoolClient, error) {
	jar := tls_client.NewCookieJar()

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Safari_IOS_16_0),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return SMSPoolClient{}, err
	}

	return SMSPoolClient{
		ApiKey:     apiKey,
		HttpClient: client,
	}, nil
}

type SMSPoolClient struct {
	ApiKey     string
	HttpClient tls_client.HttpClient
}

func (sp *SMSPoolClient) Balance() (float64, error) {
	path := `https://api.smspool.net/request/balance?key=` + sp.ApiKey

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return 0, err
	}

	resp, err := sp.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	// fmt.Println(string(body))

	return gjson.Get(string(body), "balance").Float(), nil
}

func (sp *SMSPoolClient) tryPurchaseSMS() (string, string, error) {
	path := `https://api.smspool.net/purchase/sms`

	vals := url.Values{
		"key":            {sp.ApiKey},
		"country":        {"1"},  // US
		"service":        {"39"}, // Amazon / AWS
		"pricing_option": {"1"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(vals.Encode()))
	if err != nil {
		return "", "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sp.HttpClient.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// fmt.Println(string(body))

	return gjson.Get(string(body), "phonenumber").String(), gjson.Get(string(body), "order_id").String(), nil
}

func (sp *SMSPoolClient) PurchaseSMS() (string, string, error) {
	for i := 0; i < 4; i++ {
		phone, orderid, err := sp.tryPurchaseSMS()
		if err != nil {
			continue
		}

		if phone != "" && orderid != "" {
			return phone, orderid, nil
		}

		time.Sleep(7 * time.Second)
	}

	return "", "", nil
}

func (sp *SMSPoolClient) reqRefundNumber(orderid string) ([]byte, error) {
	path := `https://api.smspool.net/sms/cancel`

	vals := url.Values{
		"key":     {sp.ApiKey},
		"orderid": {orderid},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(vals.Encode()))
	if err != nil {
		return []byte{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sp.HttpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func refundSuccess(body []byte) bool {
	return strings.Contains(string(body), "The order has been cancelled")
}

func (sp *SMSPoolClient) RefundNumber(orderid string) (bool, error) {
	for i := 0; i < 5; i++ {
		body, err := sp.reqRefundNumber(orderid)
		if err != nil {
			time.Sleep(time.Second * 7)
			continue
		}

		success := refundSuccess(body)
		if success {
			return true, nil
		}

		time.Sleep(time.Second * 7)
	}

	return false, nil
}

func (sp *SMSPoolClient) reqOTPCode(orderid string) (string, error) {
	path := `https://api.smspool.net/sms/check`

	vals := url.Values{
		"key":     {sp.ApiKey},
		"orderid": {orderid},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(vals.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := sp.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// fmt.Println(string(body))

	return gjson.Get(string(body), "sms").String(), nil
}

func (sp *SMSPoolClient) GetOTPCode(orderid string) (string, error) {
	for i := 0; i < 10; i++ {
		sms, err := sp.reqOTPCode(orderid)
		if err != nil {
			continue
		}

		if sms != "" {
			return sms, nil
		}

		time.Sleep(4 * time.Second)
	}

	return "", nil
}
