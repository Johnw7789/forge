package request

import (
	"errors"
	"net/url"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

const (
	TwoFactorAttempts = 3
)

func (t *AmazonTask) req2FAConfirm(csrfToken string) (string, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReq2FAConfirm(csrfToken)

		if err == nil && resp != "" {
			return resp, nil
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("R2FAC Failure")
}

func (t *AmazonTask) req2FAPopup() (string, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReq2FAPopup()

		if err == nil && resp != "" {
			return resp, nil
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("R2FAP Failure")
}

func (t *AmazonTask) req2FASubmit(csrfToken, sharedSecretId, code string) (string, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReq2FASubmit(csrfToken, sharedSecretId, code)

		if err == nil && resp != "" {
			return resp, nil
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("R2FAS Failure")
}

func (t *AmazonTask) req2FAPage() (string, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReq2FAPage()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302 || resp.StatusCode == 404) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("R2FAP Failure")
}

func (t *AmazonTask) reqDeletePhone(appActionToken, appAction, prevRID, workflowState string) (string, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReqDeletePhone(appActionToken, appAction, prevRID, workflowState)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("RDP Failure")
}

func (t *AmazonTask) reqPhonePage() (string, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReqPhonePage()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("RPP Failure")
}

func (t *AmazonTask) reqApproveOTP(otp, arb, returnTo, csrfToken string) (string, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReqApproveOTP(otp, arb, returnTo, csrfToken)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("RAOTP Failure")
}

func (t *AmazonTask) reqSecurityPage() (string, *http.Response, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReqSecurityPage()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			bodyStr, err := t.getBodyStr(resp)
			if err != nil {
				continue
			}

			return bodyStr, resp, nil
		}

		time.Sleep(1 * time.Second)
	}

	return "", nil, errors.New("RSP Failure")
}

func (t *AmazonTask) reqAccPage() (string, error) {
	for i := 0; i < TwoFactorAttempts; i++ {
		resp, err := t.tryReqAccPage()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302 || resp.StatusCode == 404) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("REA Failure")
}

func (t *AmazonTask) tryReq2FAConfirm(csrfToken string) (string, error) {
	path := "https://www.amazon.com/a/settings/approval/setup/enable?openid.assoc_handle=anywhere_v2_us"

	data := url.Values{
		"csrfToken": {csrfToken},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header = t.httpPostHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	return t.getBodyStr(resp)
}

func (t *AmazonTask) tryReq2FAPopup() (string, error) {
	path := "https://www.amazon.com/a/settings/approval/setup/howto?openid.assoc_handle=anywhere_v2_us"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}

	req.Header = t.httpGetSafariHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	return t.getBodyStr(resp)
}

func (t *AmazonTask) tryReq2FASubmit(csrfToken, sharedSecretId, code string) (string, error) {
	path := "https://www.amazon.com/a/settings/approval/appbackup?openid.assoc_handle=anywhere_v2_us"

	data := url.Values{
		"csrfToken":        {csrfToken},
		"isPrimary":        {"true"},
		"sharedSecretId":   {sharedSecretId},
		"sendNotification": {"false"},
		"verificationCode": {code},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header = t.httpPostHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return "", err
	}

	return t.getBodyStr(resp)
}

func (t *AmazonTask) tryReq2FAPage() (*http.Response, error) {
	path := "https://www.amazon.com/a/settings/approval/setup/register?ref_=ax_am_landing_add_2sv&openid.assoc_handle=anywhere_v2_us&openid.ns=http://specs.openid.net/auth/2.0"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = t.httpGetSafariHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqDeletePhone(appActionToken, appAction, prevRID, workflowState string) (*http.Response, error) {
	path := `https://www.amazon.com/ap/profile/mobilephone`

	data := url.Values{
		"appActionToken":    {appActionToken},
		"appAction":         {appAction},
		"prevRID":           {prevRID},
		"workflowState":     {workflowState},
		"deleteMobilePhone": {"irrelevant"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header = t.httpPostHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqPhonePage() (*http.Response, error) {
	path := "https://www.amazon.com/ap/profile/mobilephone?ref_=ax_am_landing_change_mobile&openid.assoc_handle=anywhere_v2_us&openid.ns=http://specs.openid.net/auth/2.0&referringAppAction=CNEP"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = t.httpGetSafariHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqApproveOTP(otp, arb, returnTo, csrfToken string) (*http.Response, error) {
	path := `https://www.amazon.com/ap/cvf/approval/verifyOtp`

	data := url.Values{
		"otpCode":             {otp},
		"arb":                 {arb},
		"openid.return_to":    {returnTo},
		"pageId":              {"anywhere_us"},
		"openid.assoc_handle": {"anywhere_v2_us"},
		"disableRedirect":     {"false"},
		"isResend":            {"1"},
		"csrfToken":           {csrfToken},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header = t.httpPostHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqSecurityPage() (*http.Response, error) {
	path := "https://www.amazon.com/ap/cnep?_encoding=UTF8&openid.assoc_handle=anywhere_v2_us&openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&openid.mode=checkid_setup&openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0&openid.ns.pape=http%3A%2F%2Fspecs.openid.net%2Fextensions%2Fpape%2F1.0&openid.pape.max_auth_age=0&openid.return_to=https%3A%2F%2Fwww.amazon.com%2Fgp%2Faw%2Fya%3Fie%3DUTF8%26ref_%3Dya_cnep&pageId=anywhere_us"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = t.httpGetSafariHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqAccPage() (*http.Response, error) {
	path := "https://www.amazon.com/gp/aw/ya?ref_=navm_accountmenu_account"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = t.httpGetSafariHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
