package smsman

import (
	"io"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"github.com/tidwall/gjson"
)

func NewSMSManClient(apiKey string) (SMSManClient, error) {
	jar := tls_client.NewCookieJar()

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Safari_IOS_16_0),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return SMSManClient{}, err
	}

	return SMSManClient{
		ApiKey:     apiKey,
		HttpClient: client,
	}, nil
}

type SMSManClient struct {
	ApiKey     string
	HttpClient tls_client.HttpClient
}

func (sm *SMSManClient) Balance() (float64, error) {
	path := `https://api.sms-man.com/control/get-balance?token=` + sm.ApiKey + `&currency=USD`

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return 0, err
	}

	resp, err := sm.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return gjson.Get(string(body), "balance").Float(), nil
}

func (sm *SMSManClient) tryPurchaseSMS() (string, string, error) {
	path := `https://api2.sms-man.com/control/get-number?token=` + sm.ApiKey + `&application_id=176&country_id=5&hasMultipleSms=false`

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := sm.HttpClient.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	return gjson.Get(string(body), "number").String(), gjson.Get(string(body), "request_id").String(), nil
}

func (sm *SMSManClient) PurchaseSMS() (string, string, error) {
	for i := 0; i < 4; i++ {
		phone, requestId, err := sm.tryPurchaseSMS()
		if err != nil {
			continue
		}

		if phone != "" && requestId != "" {
			return phone, requestId, nil
		}

		time.Sleep(7 * time.Second)
	}

	return "", "", nil
}

func (sm *SMSManClient) reqRefundNumber(requestId string) ([]byte, error) {
	path := `https://api2.sms-man.com/control/set-status?token=` + sm.ApiKey + `&request_id=` + requestId + `&status=reject`

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return []byte{}, err
	}

	resp, err := sm.HttpClient.Do(req)
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

func (sp *SMSManClient) RefundNumber(requestId string) (bool, error) {
	for i := 0; i < 5; i++ {
		body, err := sp.reqRefundNumber(requestId)
		if err != nil {
			time.Sleep(time.Second * 7)
			continue
		}

		success := gjson.Get(string(body), "success").Bool()
		if success {
			return true, nil
		}

		time.Sleep(time.Second * 7)
	}

	return false, nil
}

func (sm *SMSManClient) reqOTPCode(requestId string) (string, error) {
	path := `https://api2.sms-man.com/control/get-sms?token=` + sm.ApiKey + `&request_id=` + requestId

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	resp, err := sm.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return gjson.Get(string(body), "sms_code").String(), nil
}

func (sm *SMSManClient) GetOTPCode(requestId string) (string, error) {
	for i := 0; i < 10; i++ {
		sms, err := sm.reqOTPCode(requestId)
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
