package request

import (
	"errors"
	"io"
	"net/url"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

const (
	CaptchaAttempts = 3
)

func (t *AmazonTask) reqCaptchaSubmit(ci CaptchaInfo) (string, error) {
	for i := 0; i < CaptchaAttempts; i++ {
		resp, err := t.tryReqCaptchaSubmit(ci)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("Failed to Submit (RCS)")
}

func (t *AmazonTask) reqCaptchaImage(url string) ([]byte, error) {
	for i := 0; i < CaptchaAttempts; i++ {
		resp, err := t.tryReqCaptchaImage(url)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return io.ReadAll(resp.Body)
		}

		time.Sleep(1 * time.Second)
	}

	return []byte{}, errors.New("Failed to Get RCI")
}

func (t *AmazonTask) reqAddNumber(phone, verifyToken string) (string, error) {
	const CaptchaAttempts = 3

	for i := 0; i < CaptchaAttempts; i++ {
		resp, err := t.tryReqAddNumber(phone, verifyToken)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("Failed to Add Number")
}

func (t *AmazonTask) reqVerifyNumber(otp, verifyToken string) (*http.Response, error) {
	for i := 0; i < CaptchaAttempts; i++ {
		resp, err := t.tryReqVerifyNumber(otp, verifyToken)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302 || resp.StatusCode == 404) {
			return resp, nil
		}

		time.Sleep(1 * time.Second)
	}

	return nil, errors.New("Failed to Verify Number")
}

func (t *AmazonTask) tryReqChallengeJs() (*http.Response, error) {
	path := `https://ait.2608283a.us-east-1.captcha.awswaf.com/ait/ait/ait/captcha.js`

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = t.httpGetHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (t *AmazonTask) tryReqIVSitekey() (*http.Response, error) {
	path := `https://ait.2608283a.us-east-1.captcha.awswaf.com/ait/ait/ait/problem?kind=visual&domain=www.amazon.com&locale=en-us&problem=toycarcity`

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = t.httpGetHeaders()

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqCaptchaSubmit(ci CaptchaInfo) (*http.Response, error) {
	path := `https://www.amazon.com/ap/cvf/verify`

	data := url.Values{
		"cvf_captcha_captcha_token":        {ci.CaptchaToken},
		"cvf_captcha_captcha_type":         {ci.CaptchaType},
		"cvf_captcha_js_enabled_metric":    {"1"},
		"clientContext":                    {ci.ClientContext},
		"openid.pape.max_auth_age":         {"0"},
		"forceMobileLayout":                {"1"},
		"accountStatusPolicy":              {"P1"},
		"openid.identity":                  {"http://specs.openid.net/auth/2.0/identifier_select"},
		"language":                         {"en_US"},
		"pageId":                           {"amzn_device_ios_light"},
		"openid.return_to":                 {"https://www.amazon.com/ap/maplanding"},
		"openid.assoc_handle":              {"amzn_device_ios_us"},
		"openid.mode":                      {"checkid_setup"},
		"openid.ns.pape":                   {"http://specs.openid.net/extensions/pape/1.0"},
		"openid.ns.oa2":                    {"http://www.amazon.com/ap/ext/oauth/2"},
		"openid.claimed_id":                {"http://specs.openid.net/auth/2.0/identifier_select"},
		"openid.ns":                        {"http://specs.openid.net/auth/2.0"},
		"verifyToken":                      {ci.VerifyToken},
		"cvf_captcha_input":                {ci.CaptchaInput},
		"cvf_captcha_captcha_action":       {"verifyCaptcha"},
		"openid.oa2.response_type":         {"code"},
		"openid.oa2.code_challenge_method": {"S256"},
		"openid.oa2.code_challenge":        {ci.CodeChallenge},
		"openid.oa2.scope":                 {"device_auth_access"},
		"openid.oa2.client_id":             {"device:" + ci.ClientId},
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

func (t *AmazonTask) tryReqCaptchaImage(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"accept":          {"image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8"},
		"sec-fetch-site":  {"cross-site"},
		"sec-fetch-mode":  {"no-cors"},
		"sec-fetch-dest":  {"image"},
		"user-agent":      {t.taskData.UserAgent},
		"accept-language": {"en-US,en;q=0.9"},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqAddNumber(phone, verifyToken string) (*http.Response, error) {
	path := `https://www.amazon.com/ap/cvf/verify`

	data := url.Values{
		"forceMobileLayout":   {"1"},
		"openid.assoc_handle": {"amzn_device_ios_us"},
		"openid.mode":         {"checkid_setup"},
		"language":            {"en_US"},
		"openid.ns":           {"http://specs.openid.net/auth/2.0"},
		"verifyToken":         {verifyToken},
		"cvf_phone_cc":        {"US"},
		"cvf_phone_num":       {phone},
		"cvf_action":          {"collect"},
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

func (t *AmazonTask) tryReqVerifyNumber(otp, verifyToken string) (*http.Response, error) {
	path := `https://www.amazon.com/ap/cvf/verify`

	data := url.Values{
		"forceMobileLayout":   {"1"},
		"openid.assoc_handle": {"amzn_device_ios_us"},
		"openid.mode":         {"checkid_setup"},
		"language":            {"en_US"},
		"openid.ns":           {"http://specs.openid.net/auth/2.0"},
		"verifyToken":         {verifyToken},
		"code":                {otp},
		"cvf_action":          {"code"},
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
