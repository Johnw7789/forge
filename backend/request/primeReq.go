package request

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
)

const (
	PrimeAttempts = 3
)

func (t *AmazonTask) reqPrimeFinalize(pd PrimeSignupData) (string, error) {
	for i := 0; i < PrimeAttempts; i++ {
		resp, err := t.tryReqPrimeFinalize(pd)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}
	}

	return "", errors.New("Failed to Finalize Prime")
}

func (t *AmazonTask) reqPrimeWidget() (string, error) {
	for i := 0; i < PrimeAttempts; i++ {
		resp, err := t.tryReqPrimeWidget()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}
	}

	return "", errors.New("Failed to Get Prime Widget")
}

func (t *AmazonTask) reqPrimeSubmit(pd PrimeSignupData) (string, error) {
	for i := 0; i < PrimeAttempts; i++ {
		resp, err := t.tryReqPrimeSubmit(pd)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}
	}

	return "", errors.New("Failed to Submit Prime")
}

func (t *AmazonTask) reqPrimeSignup() (string, error) {
	for i := 0; i < PrimeAttempts; i++ {
		resp, err := t.tryReqPrimePage()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
			// if err != nil {
			// 	continue
			// }
			// return t.parsePrimeSignupData(bodyStr)
		}

	}

	return "", errors.New("Failed to Get Prime Signup")
}

func (t *AmazonTask) parsePrimeSignupData(resp string) (PrimeSignupData, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		return PrimeSignupData{}, errors.New("Failed to Parse PSD")
	}

	return PrimeSignupData{
		CampaignId:                 doc.Find(`input[name="primeCampaignId"]`).AttrOr("value", ""),
		OfferId:                    doc.Find(`input[name="offerId"]`).AttrOr("value", ""),
		LocationID:                 doc.Find(`input[name="locationID"]`).AttrOr("value", ""),
		OfferToken:                 doc.Find(`input[name="offerToken"]`).AttrOr("value", ""),
		RedirectURL:                doc.Find(`input[name="redirectURL"]`).AttrOr("value", ""),
		CancelRedirectURL:          doc.Find(`input[name="cancelRedirectURL"]`).AttrOr("value", ""),
		PreviousContainerRequestId: doc.Find(`input[name="previousContainerRequestId"]`).AttrOr("value", ""),
	}, nil
}

func (t *AmazonTask) parsePageId(bodyStr string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		return "", err
	}

	// Find the input element with name ppw-widgetState and get its value
	return doc.Find("input[name='actionPageDefinitionId']").AttrOr("value", ""), nil
}

func (t *AmazonTask) tryReqPrimeFinalize(pd PrimeSignupData) (*http.Response, error) {
	path := fmt.Sprintf("https://www.amazon.com/mcp/pipeline/transition?clientId=prime&locationId=wlp&offerToken=%s&session-id=%s&locationID=%s&primeCampaignId=%s&redirectURL=%s&cancelRedirectURL=%s&location=%s&paymentsPortalPreferenceType=PRIME&paymentsPortalExternalReferenceID=prime&paymentMethodId=%s&actionPageDefinitionId=%s&successUrl=/hp/wlp/pipeline/actions&failureUrl=/gp/prime/pipeline/membersignup&wlpLocation=%s&paymentMethodIdList=%s",
		pd.OfferToken,
		t.taskPrimeData.SessionId,
		pd.LocationID,
		pd.CampaignId,
		pd.RedirectURL,
		pd.CancelRedirectURL,
		pd.LocationID,
		t.taskPrimeData.PaymentMethodId,
		pd.ActionPageDefinitionId,
		pd.LocationID,
		t.taskPrimeData.PaymentMethodId,
	)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":           {"www.amazon.com"},
		"accept":         {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"sec-fetch-site": {"same-origin"},
		// "referer":         {"https://www.amazon.com/?_encoding=UTF8&ref_=navm_hdr_signin"},
		"sec-fetch-dest":  {"document"},
		"sec-fetch-mode":  {"navigate"},
		"user-agent":      {t.taskData.UserAgent},
		"accept-language": {"en-US,en;q=0.9"},
		http.HeaderOrderKey: {
			"Host",
			"accept",
			"sec-fetch-site",
			"referer",
			"sec-fetch-dest",
			"sec-fetch-mode",
			"user-agent",
			"accept-language",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqPrimeWidget() (*http.Response, error) {
	path := "https://www.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget"

	data := []KeyValue{
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.WidgetState},
		{Key: "ppw-widgetEvent", Value: "SavePaymentPreferenceEvent"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":                      {"www.amazon.com"},
		"content-type":              {"application/x-www-form-urlencoded; charset=UTF-8"},
		"widget-ajax-attempt-count": {"0"},
		"accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"x-requested-with":          {"XMLHttpRequest"},
		"sec-fetch-site":            {"same-origin"},
		"accept-language":           {"en-US,en;q=0.9"},
		"sec-fetch-mode":            {"cors"},
		"apx-widget-info":           {fmt.Sprintf("Subs:Prime/mobile/%s", t.taskPrimeData.WidgetInfo)},
		"origin":                    {"https://www.amazon.com"},
		"user-agent":                {t.taskData.UserAgent},
		// "referer":         {referer},
		"sec-fetch-dest": {"empty"},
		http.HeaderOrderKey: {
			"host",
			"content-type",
			"widget-ajax-attempt-count",
			"accept",
			"x-requested-with",
			"sec-fetch-site",
			"accept-language",
			"sec-fetch-mode",
			"apx-widget-info",
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

func (t *AmazonTask) tryReqPrimeSubmit(pd PrimeSignupData) (*http.Response, error) {
	path := "https://www.amazon.com/gp/prime/pipeline/confirm?offerToken=" + pd.OfferToken

	data := []KeyValue{
		{Key: "primeCampaignId", Value: pd.CampaignId},
		{Key: "offerId", Value: pd.OfferId},
		{Key: "locationID", Value: pd.LocationID},
		{Key: "offerToken", Value: pd.OfferToken},
		{Key: "redirectURL", Value: pd.RedirectURL},
		{Key: "cancelRedirectURL", Value: pd.CancelRedirectURL},
		{Key: "previousContainerRequestId", Value: pd.PreviousContainerRequestId},
		{Key: "CTAtext", Value: "Submit"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
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
		// "referer":        {"https://www.amazon.com/amazonprime?ref_=navm_em_allpf_prime_nonmember_0_1_1_67"},
		"sec-fetch-dest": {"document"},
		http.HeaderOrderKey: {
			"Host",
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

func (t *AmazonTask) tryReqPrimePage() (*http.Response, error) {
	path := "https://www.amazon.com/mc?_encoding=UTF8&ref_=navm_accountmenu_prime"

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
