package request

import (
	"errors"
	"strings"
	"time"
)

func (t *AmazonTask) submitSMS(provider, verifyToken string) error {
	switch provider {
	case "SMS Pool":
		err := t.handleSPFlow(verifyToken)
		if err != nil {
			return err
		}
	case "SMS Man":
		err := t.handleSMFlow(verifyToken)
		if err != nil {
			return err
		}
	case "Daisy SMS":
		err := t.handleDaisyFlow(verifyToken)
		if err != nil {
			return err
		}
	default:
		return errors.New("Invalid SMS Provider")
	}

	return nil
}

func (t *AmazonTask) handleSPFlow(verifyToken string) error {
	if t.spClient.ApiKey == "" {
		return errors.New("SMS Pool Client is not initialized")
	}

	var bodyStr string
	var err error

	var number string
	var orderid string

	for i := 0; i < t.UserInfo.SmsInfo.MaxRetries; i++ {
		t.UpdateStatus("Retrieving Number")
		number, orderid, err = t.spClient.PurchaseSMS()
		if err != nil {
			continue
		}

		if number == "" || orderid == "" {
			continue
		}

		t.UpdateStatus("Retrieved SMS Number: " + number)

		bodyStr, err = t.reqAddNumber(number, verifyToken)
		if err != nil {
			t.UpdateStatus("Refunding Number: " + number)
			success, err := t.spClient.RefundNumber(orderid)
			if err != nil || !success {
				continue
			}

			t.UpdateStatus("Refunded Number: " + number)
			time.Sleep(time.Second * 3)

			continue
		}

		newToken, _ := t.evalCreateAccount(bodyStr)
		if newToken != "" {
			verifyToken = newToken
		}

		code, err := t.spClient.GetOTPCode(orderid)
		if err != nil || code == "" {
			t.UpdateStatus("Refunding Number: " + number)
			success, err := t.spClient.RefundNumber(orderid)
			if err != nil || !success {
				continue
			}

			t.UpdateStatus("Refunded Number: " + number)
			time.Sleep(time.Second * 3)

			continue
		}

		t.UpdateStatus("Received SMS Code: " + code)

		resp, err := t.reqVerifyNumber(code, verifyToken)
		if err != nil {
			t.UpdateStatus("Retrying SMS")
			time.Sleep(time.Second * 3)
			continue
		}

		bodyStr, err = t.getBodyStr(resp)
		if err != nil {
			t.UpdateStatus("Retrying SMS")
			time.Sleep(time.Second * 3)
			continue
		}

		if !strings.Contains(bodyStr, "Add mobile number") {
			// lcoation := resp.Request.URL.String()
			// t.taskData.AuthCode, _ = t.getAuthorizationCode(lcoation)

			break
		}
	}

	if strings.Contains(bodyStr, "Add mobile number") {
		return errors.New("Unable to Verify SMS")
	}

	if err != nil {
		return errors.New("Error - Unable to Verify SMS")
	}

	t.Phone = number
	return nil
}

func (t *AmazonTask) handleSMFlow(verifyToken string) error {
	if t.smClient.ApiKey == "" {
		return errors.New("SMS Man Client is not initialized")
	}

	var bodyStr string
	var err error

	var number string
	var requestId string

	for i := 0; i < t.UserInfo.SmsInfo.MaxRetries; i++ {
		t.UpdateStatus("Retrieving Number")
		number, requestId, err = t.smClient.PurchaseSMS()
		if err != nil {
			continue
		}

		if number == "" || requestId == "" {
			continue
		}

		t.UpdateStatus("Retrieved SMS Number: " + number)

		bodyStr, err = t.reqAddNumber(number, verifyToken)
		if err != nil {
			t.UpdateStatus("Refunding Number: " + number)
			success, err := t.smClient.RefundNumber(requestId)
			if err != nil || !success {
				t.UpdateStatus("Retrying SMS")
				time.Sleep(time.Second * 3)
				continue
			}

			t.UpdateStatus("Refunded Number: " + number)
			time.Sleep(time.Second * 3)

			continue
		}

		newToken, _ := t.evalCreateAccount(bodyStr)
		if newToken != "" {
			verifyToken = newToken
		}

		code, err := t.smClient.GetOTPCode(requestId)
		if err != nil || code == "" {
			t.UpdateStatus("Refunding Number: " + number)
			success, err := t.smClient.RefundNumber(requestId)
			if err != nil || !success {
				t.UpdateStatus("Retrying SMS")
				time.Sleep(time.Second * 3)
				continue
			}

			t.UpdateStatus("Refunded Number: " + number)
			time.Sleep(time.Second * 3)

			continue
		}

		t.UpdateStatus("Received SMS Code: " + code)

		resp, err := t.reqVerifyNumber(code, verifyToken)
		if err != nil {
			t.UpdateStatus("Retrying SMS")
			time.Sleep(time.Second * 3)
			continue
		}

		bodyStr, err = t.getBodyStr(resp)
		if err != nil {
			t.UpdateStatus("Retrying SMS")
			time.Sleep(time.Second * 3)
			continue
		}

		if !strings.Contains(bodyStr, "Add mobile number") {
			// lcoation := resp.Request.URL.String()
			// t.taskData.AuthCode, _ = t.getAuthorizationCode(lcoation)

			break
		}
	}

	if strings.Contains(bodyStr, "Add mobile number") {
		return errors.New("Unable to Verify SMS")
	}

	if err != nil {
		return errors.New("Error - Unable to Verify SMS")
	}

	t.Phone = number

	return nil
}

func (t *AmazonTask) handleDaisyFlow(verifyToken string) error {
	if t.dsClient.ApiKey == "" {
		return errors.New("Daisy SMS Client is not initialized")
	}

	var bodyStr string
	var err error

	var number string
	var orderId string

	for i := 0; i < t.UserInfo.SmsInfo.MaxRetries; i++ {
		t.UpdateStatus("Retrieving Number")
		orderId, number, err = t.dsClient.PurchaseSMS()
		if err != nil {
			continue
		}

		if orderId == "" || number == "" {
			continue
		}

		t.UpdateStatus("Retrieved SMS Number: " + number)

		bodyStr, err = t.reqAddNumber(number, verifyToken)
		if err != nil {
			t.UpdateStatus("Refunding Number: " + number)
			success, err := t.dsClient.RefundNumber(orderId)
			if err != nil || !success {
				t.UpdateStatus("Retrying SMS")
				time.Sleep(time.Second * 3)
				continue
			}

			t.UpdateStatus("Refunded Number: " + number)
			time.Sleep(time.Second * 3)

			continue
		}

		newToken, _ := t.evalCreateAccount(bodyStr)
		if newToken != "" {
			verifyToken = newToken
		}

		code, err := t.dsClient.GetOTPCode(orderId)
		if err != nil || code == "" {
			t.UpdateStatus("Refunding Number: " + number)
			success, err := t.dsClient.RefundNumber(orderId)
			if err != nil || !success {
				t.UpdateStatus("Retrying SMS")
				time.Sleep(time.Second * 3)
				continue
			}

			t.UpdateStatus("Refunded Number: " + number)
			time.Sleep(time.Second * 3)

			continue
		}

		t.UpdateStatus("Received SMS Code: " + code)

		resp, err := t.reqVerifyNumber(code, verifyToken)
		if err != nil {
			t.UpdateStatus("Retrying SMS")
			time.Sleep(time.Second * 3)
			continue
		}

		bodyStr, err = t.getBodyStr(resp)
		if err != nil {
			t.UpdateStatus("Retrying SMS")
			time.Sleep(time.Second * 3)
			continue
		}

		if !strings.Contains(bodyStr, "Add mobile number") {
			// lcoation := resp.Request.URL.String()
			// t.taskData.AuthCode, _ = t.getAuthorizationCode(lcoation)

			break
		}
	}

	if strings.Contains(bodyStr, "Add mobile number") {
		return errors.New("Unable to Verify SMS")
	}

	if err != nil {
		return errors.New("Error - Unable to Verify SMS")
	}

	t.Phone = number

	return nil
}
