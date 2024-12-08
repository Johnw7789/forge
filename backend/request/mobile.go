package request

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

// *** Alternate mobile app flow to the Alexa flow ***

func (t *AmazonTask) reqVerifyEmailMobile(code, verifyToken string) (string, error) {
	const attempts = 3

	for i := 0; i < attempts; i++ {
		resp, err := t.tryReqVerifyEmailMobile(code, verifyToken)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("Failed to Verify Email")
}

func (t *AmazonTask) reqCreateAccountMobile(sessionId string) (string, error) {
	const attempts = 3

	for i := 0; i < attempts; i++ {
		resp, err := t.tryReqCreateAccountMobile(sessionId)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)
	}

	return "", errors.New("Failed to Create Account")
}

func (t *AmazonTask) reqCreateDataMobile() (string, error) {
	const attempts = 3

	for i := 0; i < attempts; i++ {
		resp, err := t.tryReqCreateDataMobile()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)

	}

	return "", errors.New("Failed to Get RCDI")
}

func (t *AmazonTask) reqMShop() (string, error) {
	const attempts = 3

	for i := 0; i < attempts; i++ {
		resp, err := t.tryReqMShop()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)

	}

	return "", errors.New("Failed to Get Orders")
}

func (t *AmazonTask) reqHomePage() (string, error) {
	const attempts = 3

	for i := 0; i < attempts; i++ {
		resp, err := t.tryReqHomePage()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(1 * time.Second)

	}

	return "", errors.New("Failed to Get RHPI")
}

func (t *AmazonTask) tryReqVerifyEmailMobile(code, verifyToken string) (*http.Response, error) {
	path := `https://www.amazon.com/ap/cvf/verify`

	data := url.Values{
		"action":              {"code"},
		"forceMobileLayout":   {"1"},
		"openid.assoc_handle": {"amzn_device_ios_us"},
		"openid.mode":         {"checkid_setup"},
		"openid.ns":           {"http://specs.openid.net/auth/2.0"},
		"verifyToken":         {verifyToken},
		"code":                {code},
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

func (t *AmazonTask) tryReqCreateAccountMobile(sessionId string) (*http.Response, error) {
	path := `https://www.amazon.com/ap/register`
	if sessionId != "" {
		path = path + "/" + sessionId
	}

	data := url.Values{
		"appActionToken":      {t.taskData.CreateInfo.AppActionToken},
		"appAction":           {t.taskData.CreateInfo.AppAction},
		"openid.return_to":    {t.taskData.CreateInfo.OpenIdReturnTo},
		"prevRID":             {t.taskData.CreateInfo.PrevRID},
		"workflowState":       {t.taskData.CreateInfo.WorkflowState},
		"customerName":        {t.UserInfo.FirstName + " " + t.UserInfo.LastName},
		"email":               {t.UserInfo.Email},
		"password":            {t.UserInfo.Password},
		"showPasswordChecked": {"true"},
		// "metadata1":           {t.taskInfo.Metadata},
	}

	// if sessionId == "" {
	// 	data.Add("metadata1", t.taskInfo.Metadata)
	// }

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

func (t *AmazonTask) tryReqLocationArb(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
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

func (t *AmazonTask) tryReqCreateDataMobile() (*http.Response, error) {
	path := `https://www.amazon.com/ap/register?openid.return_to=https://www.amazon.com&openid.assoc_handle=amzn_device_ios_us&openid.identity=http://specs.openid.net/auth/2.0/identifier_select&pageId=amzn_device_ios_light&accountStatusPolicy=P1&openid.claimed_id=http://specs.openid.net/auth/2.0/identifier_select&openid.mode=checkid_setup&openid.ns.oa2=http://www.amazon.com/ap/ext/oauth/2&language=en_US&openid.ns.pape=http://specs.openid.net/extensions/pape/1.0&openid.ns=http://specs.openid.net/auth/2.0&openid.pape.max_auth_age=0`

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

func (t *AmazonTask) tryReqMShop() (*http.Response, error) {
	req, err := http.NewRequest("GET", "https://www.amazon.com/gp/gw/ajax/mshop.html?modern=1", nil)
	if err != nil {
		return nil, err
	}

	// 'Host': 'www.amazon.com',
	// # 'Cookie': 'amzn-app-ctxt=1.8%20{"xv":"1.15"%2C"ov":"17.2.1"%2C"uiv":5%2C"an":"Amazon"%2C"cp":986000%2C"os":"iOS"%2C"dm":{"pb":"94"%2C"w":"1290"%2C"pt":"69"%2C"h":"2503"%2C"ld":"3.000000"}%2C"msd":".amazon.com"%2C"di":{"ca":"--"%2C"dsn":"EF23B177-8B7E-4094-BBE6-6B7F53FA967A"%2C"mf":"Apple"%2C"ct":"Wifi"%2C"pr":"iPhone"%2C"md":"iPhone"%2C"v":"Unknown_iPhone15%2C3"%2C"dti":"A287KHUN77EJVL"}%2C"ast":3%2C"aid":"com.amazon.Amazon"%2C"av":"23.7.2"}; amzn-app-id=Amazon/23.7.2/1-634215.0; mobile-device-info=scale:3|w:430|h:834; lc-main=en_US; session-token=daistjZK0a2kBB1xTzspQjrRpEXNX7ESbjYmXFfylVMu83GZaeWpiCMYl09xT1QP6ZDhfRC1cMOSjWyPOQvTRTqjKLBTKSC5T8uedSAqiEpYug2BO4AcKiJBKdjsCVqRb1nrPcp14MEnDkrSuVpOg/kOLBU7uPGdVCMhm8mADMd+K06xBL9xe89RZ7HJ6udxflMkWsvZUFoP2spqWDBvUCqnKT5C9jvkIdcdjuAWSm1cYUoUatuTKKq4rdTJS5H6GclOLaKyG3cbqFZJJCtUK7BdN8i/+TiENUzbXtfjTyLYXED09EoKnjB7CZlxOkT+AXm9WLF7h+M2FY6i13+0fGDni7hCipORYnII320qdnVKSGY0KZLF6bN7O0jz/UIS; csm-hit=NF0FKB8AEAYG9EZKG10T+s-03R1PEZ8GQD40ESEY1GG|1711690432002; at-main="Atza|IwEBIGRvZII8dkq_ZMIQYQrnq-Nxzokj4GhPiB7ooTuHTDPhMxpxerw0ETzNqW2a3YnexEhWgHuJ17UgcsjgJg7waBiqZP26y91Q3Pb-LtporvjUkAN0l_b-pFLTCfJhKFWkTkXBFfa2pdw19sYWswdtx6-lpuHIMwTMNelpjZyzAS11MK9D77mbLPo_xhoAKaPXwwfexux4S6g524G9I_p2S-vie-nnq_4vTIkAGS9ZSc8TLxflNaAe0vPJWKFa8Ss4jZB4BKeCtyHm9REOa5BieJ-IPVPP5z3d_1pcBxGAu9xrOLtGbWB-4i9uQ4jNxVzY8j0"; sess-at-main="zv1qggqTpgbYNuEpgfQDy/28A/ZPRTUHvfs3FMUG31s="; session-id=133-8535439-7108425; ubid-main=132-3615163-7343656; x-main="TA@JnTE0FynrzqA6X?AZSKEKFo?sWKMsz5k2hKcOp3ZIIczxMtwlD5XBI@XR1vHr"; i18n-prefs=USD; privacy-consent=%7B%22avlString%22%3A%22%22%2C%22gvlString%22%3A%22%22%2C%22amazonAdvertisingPublisher%22%3Atrue%7D; session-id-time=2082787201l',
	// 'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8',
	// 'sec-fetch-site': 'none',
	// 'x-mash-csm-tc': '1711690440161',
	// 'sec-fetch-mode': 'navigate',
	// 'user-agent': 'Mozilla/5.0 (iPhone; CPU iPhone OS 17_2_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148',
	// 'accept-language': 'en-US,en;q=0.9',
	// 'sec-fetch-dest': 'document',

	req.Header = http.Header{
		"Host":            {"www.amazon.com"},
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"sec-fetch-site":  {"none"},
		"x-mash-csm-tc":   {fmt.Sprintf("%d", time.Now().UnixMilli())},
		"sec-fetch-mode":  {"navigate"},
		"user-agent":      {t.taskData.UserAgent},
		"accept-language": {"en-US,en;q=0.9"},
		"sec-fetch-dest":  {"document"},
		http.HeaderOrderKey: {
			"Host",
			"cookie",
			"accept",
			"sec-fetch-site",
			"x-mash-csm-tc",
			"sec-fetch-mode",
			"user-agent",
			"accept-language",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqHomePage() (*http.Response, error) {
	req, err := http.NewRequest("GET", "https://www.amazon.com/", nil)
	if err != nil {
		return nil, err
	}

	req.Header = t.httpGetHeaders()

	// req.Header.Add("referer", "https://www.google.com/")

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
