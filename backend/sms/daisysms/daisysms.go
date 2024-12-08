package daisysms

import (
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

type DaisySMSClient struct {
	ApiKey     string
	HttpClient tls_client.HttpClient
}

func NewDaisySMSClient(apiKey string) (DaisySMSClient, error) {
	jar := tls_client.NewCookieJar()

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(profiles.Safari_IOS_16_0),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return DaisySMSClient{}, err
	}

	return DaisySMSClient{
		ApiKey:     apiKey,
		HttpClient: client,
	}, nil
}

func (ds *DaisySMSClient) Balance() (float64, error) {
	path := `https://daisysms.com/stubs/handler_api.php?api_key=` + ds.ApiKey + `&action=getBalance`

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return 0, err
	}

	resp, err := ds.HttpClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	balanceSpl := strings.Split(string(body), ":")
	if len(balanceSpl) < 2 || balanceSpl[0] != "ACCESS_BALANCE" {
		return 0, errors.New("Failed to parse balance: " + string(body))
	}

	// convert balanceSpl[1] to float64
	balance, err := strconv.ParseFloat(balanceSpl[1], 64)
	if err != nil {
		return 0, errors.New("Failed to parse balance: " + string(body))
	}

	return balance, nil
}

func (ds *DaisySMSClient) PurchaseSMS() (string, string, error) {
	for i := 0; i < 4; i++ {
		orderId, phone, err := ds.tryPurchaseSMS()
		if err != nil {
			continue
		}

		if orderId != "" && phone != "" {
			// trim 1 at beginning of phone number
			phone = phone[1:]
			return orderId, phone, nil
		}

		time.Sleep(7 * time.Second)
	}

	return "", "", nil
}

func (ds *DaisySMSClient) RefundNumber(requestId string) (bool, error) {
	for i := 0; i < 5; i++ {
		body, err := ds.reqRefundNumber(requestId)
		if err != nil {
			time.Sleep(time.Second * 7)
			continue
		}

		success := strings.Contains(string(body), "ACCESS_CANCEL")
		if success {
			return true, nil
		}

		time.Sleep(time.Second * 7)
	}

	return false, nil
}

func (ds *DaisySMSClient) GetOTPCode(requestId string) (string, error) {
	for i := 0; i < 10; i++ {
		sms, err := ds.reqOTPCode(requestId)
		if err != nil {
			time.Sleep(4 * time.Second)
			continue
		}

		if sms != "" {
			return sms, nil
		}

		time.Sleep(4 * time.Second)
	}

	return "", nil
}

func (ds *DaisySMSClient) tryPurchaseSMS() (string, string, error) {
	path := `https://daisysms.com/stubs/handler_api.php?api_key=` + ds.ApiKey + `&action=getNumber&service=am&max_price=0.5`

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := ds.HttpClient.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	numberSpl := strings.Split(string(body), ":")
	if len(numberSpl) < 3 || numberSpl[0] != "ACCESS_NUMBER" {
		return "", "", errors.New("Failed to Retrieve Number")
	}

	return numberSpl[1], numberSpl[2], nil
}

func (ds *DaisySMSClient) reqRefundNumber(orderId string) ([]byte, error) {
	path := `https://daisysms.com/stubs/handler_api.php?api_key=` + ds.ApiKey + `&action=setStatus&id=` + orderId + `&status=8`

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return []byte{}, err
	}

	resp, err := ds.HttpClient.Do(req)
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

func (ds *DaisySMSClient) reqOTPCode(orderId string) (string, error) {
	path := `https://daisysms.com/stubs/handler_api.php?api_key=` + ds.ApiKey + `&action=getStatus&id=` + orderId

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	resp, err := ds.HttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	codeSpl := strings.Split(string(body), ":")
	if len(codeSpl) < 2 || !strings.Contains(codeSpl[0], "STATUS_OK") {
		return "", nil
	}

	return codeSpl[1], nil
}
