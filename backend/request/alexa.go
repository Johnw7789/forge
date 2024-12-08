package request

import (
	"errors"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

const (
	CreateAttemtps = 3
)

func (t *AmazonTask) reqCreateAccountAlexa() (string, error) {
	for i := 0; i < CreateAttemtps; i++ {
		resp, err := t.tryReqCreateAccountAlexa()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Create Account")
}

func (t *AmazonTask) reqCreateDataAlexa(url, referer string) (string, error) {
	for i := 0; i < CreateAttemtps; i++ {
		resp, err := t.tryReqCreateDataAlexa(url, referer)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)

	}

	return "", errors.New("Failed to Get Create Data")
}

func (t *AmazonTask) reqSigninAlexa() (string, string, error) {
	for i := 0; i < CreateAttemtps; i++ {
		resp, err := t.tryReqSigninAlexa()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			finalUrl := resp.Request.URL.String()
			bodyStr, err := t.getBodyStr(resp)
			return bodyStr, finalUrl, err
		}

		time.Sleep(2 * time.Second)

	}

	return "", "", errors.New("Failed to Get Signin")
}

func (t *AmazonTask) tryReqCreateAccountAlexa() (*http.Response, error) {
	path := `https://www.amazon.com/ap/register`

	data := []KeyValue{
		{"appActionToken", t.taskData.CreateInfo.AppActionToken},
		{"appAction", t.taskData.CreateInfo.AppAction},
		{"openid.return_to", t.taskData.CreateInfo.OpenIdReturnTo},
		{"prevRID", t.taskData.CreateInfo.PrevRID},
		{"workflowState", t.taskData.CreateInfo.WorkflowState},
		{"customerName", t.UserInfo.FirstName + " " + t.UserInfo.LastName},
		{"email", t.UserInfo.Email},
		{"password", t.UserInfo.Password},
		{"showPasswordChecked", "true"},
		// {"metadata1", t.taskData.Metadata},
	}

	encodedData := encodeData(data)

	req, err := http.NewRequest("POST", path, strings.NewReader(encodedData))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":            {"www.amazon.com"},
		"content-type":    {"application/x-www-form-urlencoded"},
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"sec-fetch-site":  {"same-origin"},
		"accept-language": {"en-US,en;q=0.9"},
		"sec-fetch-mode":  {"navigate"},
		"origin":          {"https://www.amazon.com"},
		"user-agent":      {t.taskData.UserAgent},
		"sec-fetch-dest":  {"document"},
		http.HeaderOrderKey: {
			"Host",
			"Cookie",
			"content-type",
			"accept",
			"sec-fetch-site",
			"accept-language",
			"sec-fetch-mode",
			"origin",
			"user-agent",
			"referer",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqCreateDataAlexa(url, referer string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":            {"www.amazon.com"},
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"sec-fetch-site":  {"none"},
		"sec-fetch-dest":  {"document"},
		"accept-language": {"en-US,en;q=0.9"},
		"sec-fetch-mode":  {"navigate"},
		"user-agent":      {t.taskData.UserAgent},
		"referer":         {referer},
		http.HeaderOrderKey: {
			"Host",
			"Cookie",
			"accept",
			"sec-fetch-site",
			"sec-fetch-dest",
			"accept-language",
			"sec-fetch-mode",
			"user-agent",
			"referer",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqSigninAlexa() (*http.Response, error) {
	urlBase := "https://www.amazon.com/ap/signin"

	params := []KeyValue{
		{"openid.return_to", "https://www.amazon.com/ap/maplanding"},
		{"openid.oa2.code_challenge_method", "S256"},
		{"openid.assoc_handle", "amzn_dp_project_dee_ios"},
		{"openid.identity", "http://specs.openid.net/auth/2.0/identifier_select"},
		{"pageId", "amzn_dp_project_dee_ios"},
		{"accountStatusPolicy", "P1"},
		{"openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select"},
		{"openid.mode", "checkid_setup"},
		{"openid.ns.oa2", "http://www.amazon.com/ap/ext/oauth/2"},
		{"openid.oa2.client_id", "device:" + t.taskData.ChallengeData.ClientId},
		{"language", "en_US"},
		{"openid.ns.pape", "http://specs.openid.net/extensions/pape/1.0"},
		{"openid.oa2.code_challenge", t.taskData.ChallengeData.VerifierChecksum},
		{"openid.oa2.scope", "device_auth_access"},
		{"openid.ns", "http://specs.openid.net/auth/2.0"},
		{"openid.pape.max_auth_age", "0"},
		{"openid.oa2.response_type", "code"},
	}

	path := "?" + encodeData(params)

	req, err := http.NewRequest("GET", urlBase+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":            {"www.amazon.com"},
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"accept-charset":  {"utf-8"},
		"sec-fetch-site":  {"none"},
		"accept-language": {"en-US"},
		"cache-control":   {"no-store"},
		"sec-fetch-mode":  {"navigate"},
		"user-agent":      {t.taskData.UserAgent},
		"sec-fetch-dest":  {"document"},
		http.HeaderOrderKey: {
			"Host",
			"Cookie",
			"accept",
			"accept-charset",
			"sec-fetch-site",
			"accept-language",
			"cache-control",
			"sec-fetch-mode",
			"user-agent",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
