package request

import (
	"errors"
	"fmt"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
)

const (
	AmazonPrimeURL = "https://www.amazon.com/amazonprime"
	InfoAttempts   = 3
)

func (t *AmazonTask) reqFinalize3() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqFinalize3()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Finalize (3)")
}

func (t *AmazonTask) reqFinalize2() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqFinalize2()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Finalize (2)")
}

func (t *AmazonTask) reqFinalize1() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqFinalize1()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Finalize")
}

func (t *AmazonTask) reqAddressSelect() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqAddressSelect()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Select Address")
}

func (t *AmazonTask) reqAddressId() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqAddressId()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Get AID")
}

func (t *AmazonTask) reqSubmitPayment(cardType string) (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqSubmitPayment(cardType)

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Submit Payment")
}

func (t *AmazonTask) reqCardType() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqCardType()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Get Card Type")
}

func (t *AmazonTask) reqRegisterWidget() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqRegisterWidget()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Register Widget")
}

func (t *AmazonTask) reqAddAddress() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqAddAddress()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Add Address")
}

func (t *AmazonTask) reqAddPurchasePref() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReqAddPurchasePref()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Add Purchase Pref")
}

func (t *AmazonTask) req1ClickPage() (string, error) {
	for i := 0; i < InfoAttempts; i++ {
		resp, err := t.tryReq1ClickPage()

		if err == nil && resp != nil && (resp.StatusCode == 200 || resp.StatusCode == 302) {
			return t.getBodyStr(resp)
		}

		time.Sleep(2 * time.Second)
	}

	return "", errors.New("Failed to Get ACM")
}

func (t *AmazonTask) tryReqFinalize3() (*http.Response, error) {
	path := "https://www.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget"

	data := []KeyValue{
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.FinalWidgetState},
		{Key: "ppw-widgetEvent", Value: "SavePaymentPreferenceEvent"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":                      {"apx-security.amazon.com"},
		"Widget-Ajax-Attempt-Count": {"0"},
		"Accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"X-Requested-With":          {"XMLHttpRequest"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Sec-Fetch-Mode":            {"cors"},
		"APX-Widget-Info":           {fmt.Sprintf("YA:OneClick/mobile/%s", t.taskPrimeData.WidgetInfo)},
		"Origin":                    {"https://apx-security.amazon.com"},
		"User-Agent":                {t.taskData.UserAgent},
		"Content-Type":              {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Sec-Fetch-Dest":            {"empty"},
		http.HeaderOrderKey: {
			"host",
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
			"content-type",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqFinalize2() (*http.Response, error) {
	path := "https://www.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget"

	data := []KeyValue{
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.FinalWidgetState},
		{Key: "ie", Value: "UTF-8"},
		{Key: "ppw-instrumentRowSelection", Value: "instrumentId=" + t.taskPrimeData.InstrumentId + "&isExpired=false&paymentMethod=CC&tfxEligible=false"},
		{Key: "ppw-" + t.taskPrimeData.InstrumentId + "_instrumentOrderTotalBalance", Value: "{}"},
		{Key: "ppw-widgetEvent:PreferencePaymentOptionSelectionEvent", Value: ""},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":                      {"apx-security.amazon.com"},
		"Widget-Ajax-Attempt-Count": {"0"},
		"Accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"X-Requested-With":          {"XMLHttpRequest"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Sec-Fetch-Mode":            {"cors"},
		"APX-Widget-Info":           {fmt.Sprintf("YA:OneClick/mobile/%s", t.taskPrimeData.WidgetInfo)},
		"Origin":                    {"https://apx-security.amazon.com"},
		"User-Agent":                {t.taskData.UserAgent},
		"Content-Type":              {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Sec-Fetch-Dest":            {"empty"},
		http.HeaderOrderKey: {
			"host",
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
			"content-type",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqFinalize1() (*http.Response, error) {
	path := "https://www.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget"

	data := []KeyValue{
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.FinalWidgetState},
		{Key: "ppw-widgetEvent", Value: "AddPaymentMethodRefreshEvent"},
		{Key: "ppw-paymentMethodId", Value: t.taskPrimeData.InstrumentId},
		{Key: "ppw-widgetAction", Value: "add-credit-card-workflow-complete"},
		{Key: "ppw-maybeShouldRecordAPX3Metric", Value: "0"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":                      {"apx-security.amazon.com"},
		"Widget-Ajax-Attempt-Count": {"0"},
		"Accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"X-Requested-With":          {"XMLHttpRequest"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Sec-Fetch-Mode":            {"cors"},
		"APX-Widget-Info":           {fmt.Sprintf("YA:OneClick/mobile/%s", t.taskPrimeData.WidgetInfo)},
		"Origin":                    {"https://apx-security.amazon.com"},
		"User-Agent":                {t.taskData.UserAgent},
		"Content-Type":              {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Sec-Fetch-Dest":            {"empty"},
		http.HeaderOrderKey: {
			"host",
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
			"content-type",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqAddressSelect() (*http.Response, error) {
	path := "https://apx-security.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget"

	data := []KeyValue{
		{Key: `ppw-widgetEvent:SelectAddressEvent:{"addressId":"` + t.taskPrimeData.AddressId + `"}`, Value: ""},
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.WidgetState},
		{Key: "ie", Value: "UTF-8"},
		{Key: "ppw-pickAddressType", Value: "Inline"},
		{Key: "ppw-addressSelection", Value: t.taskPrimeData.AddressId},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":                      {"apx-security.amazon.com"},
		"Widget-Ajax-Attempt-Count": {"0"},
		"Accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"X-Requested-With":          {"XMLHttpRequest"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Sec-Fetch-Mode":            {"cors"},
		"APX-Widget-Info":           {fmt.Sprintf("YA:OneClick/mobile/%s", t.taskPrimeData.WidgetInfo)},
		"Origin":                    {"https://apx-security.amazon.com"},
		"User-Agent":                {t.taskData.UserAgent},
		"Content-Type":              {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Sec-Fetch-Dest":            {"empty"},
		http.HeaderOrderKey: {
			"host",
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
			"content-type",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqAddressId() (*http.Response, error) {
	path := "https://www.amazon.com/portal-migration/hz/glow/get-location-label?pageType=unknown&deviceType=mobile&osType=ios"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":            {"www.amazon.com"},
		"accept":          {"*/*"},
		"user-agent":      {t.taskData.UserAgent},
		"accept-language": {"en-US,en;q=0.9"},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqSubmitPayment(cardType string) (*http.Response, error) {
	path := "https://apx-security.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget?sif_profile=APX-Encrypt-All-NA"

	data := []KeyValue{
		{Key: "ppw-widgetEvent:AddCreditCardEvent", Value: ""},
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.WidgetState},
		{Key: "ie", Value: "UTF-8"},
		{Key: "addCreditCardNumber", Value: t.Payment.CardNumber},                                            // separate by spaces? eg visa 4 groups of 4
		{Key: "ppw-expirationDate_combinedMonthYear", Value: t.Payment.ExpMonth + " / " + t.Payment.ExpYear}, // todo: verify?? ex 09 / 24
		{Key: "ppw-accountHolderName", Value: t.Payment.CardHolder},
		{Key: "addCreditCardVerificationNumber", Value: ""},
		{Key: "ppw-addCreditCardVerificationNumber_isRequired", Value: "false"},
		{Key: "ppw-addCreditCardPostalCode", Value: ""},
		{Key: "ppw-addCreditCardPostalCode_isRequired", Value: "false"},
		{Key: "__sif_encrypted_hba_account_holder_name", Value: t.Payment.CardHolder},
		{Key: "ppw-issuer", Value: cardType},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":                      {"apx-security.amazon.com"},
		"Widget-Ajax-Attempt-Count": {"0"},
		"Accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"X-Requested-With":          {"XMLHttpRequest"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Sec-Fetch-Mode":            {"cors"},
		"APX-Widget-Info":           {fmt.Sprintf("YA:OneClick/mobile/%s", t.taskPrimeData.WidgetInfo)},
		"Origin":                    {"https://apx-security.amazon.com"},
		"User-Agent":                {t.taskData.UserAgent},
		"Content-Type":              {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Sec-Fetch-Dest":            {"empty"},
		http.HeaderOrderKey: {
			"host",
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
			"content-type",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqCardType() (*http.Response, error) {
	path := "https://apx-security.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget"

	data := []KeyValue{
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "addCreditCardNumber", Value: t.Payment.CardNumber},
		{Key: "ppw-widgetEvent", Value: "IdentifyCreditCardEvent"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.WidgetState},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":                      {"apx-security.amazon.com"},
		"Widget-Ajax-Attempt-Count": {"0"},
		"Accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"X-Requested-With":          {"XMLHttpRequest"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Accept-Language":           {"en-US,en;q=0.9"},
		"Sec-Fetch-Mode":            {"cors"},
		"APX-Widget-Info":           {fmt.Sprintf("YA:OneClick/mobile/%s", t.taskPrimeData.WidgetInfo)},
		"Origin":                    {"https://apx-security.amazon.com"},
		"User-Agent":                {t.taskData.UserAgent},
		"Content-Type":              {"application/x-www-form-urlencoded; charset=UTF-8"},
		"Sec-Fetch-Dest":            {"empty"},
		http.HeaderOrderKey: {
			"host",
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
			"content-type",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqAddAddress() (*http.Response, error) {
	path := "https://www.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget"

	data := []KeyValue{
		{Key: "ppw-widgetEvent:AddAddressEvent", Value: ""},
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.WidgetState},
		{Key: "ie", Value: "UTF-8"},
		{Key: "ppw-fullName", Value: t.Address.FullName},
		{Key: "ppw-line1", Value: t.Address.Line1},
		{Key: "ppw-line2", Value: t.Address.Line2},
		{Key: "ppw-city", Value: t.Address.City},
		{Key: "ppw-stateOrRegion", Value: t.Address.State},
		{Key: "ppw-postalCode", Value: t.Address.Zip},
		{Key: "ppw-phoneNumber", Value: t.Address.PhoneNumber},
		{Key: "ppw-countryCode", Value: "US"},
		{Key: "ppw-enableOneClick", Value: "enableOneClick"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":                      {"www.amazon.com"},
		"Widget-Ajax-Attempt-Count": {"0"},
		"accept":                    {"application/json, text/javascript, */*; q=0.01"},
		"x-requested-with":          {"XMLHttpRequest"},
		"sec-fetch-site":            {"same-origin"},
		"accept-language":           {"en-US,en;q=0.9"},
		"sec-fetch-mode":            {"cors"},
		"apx-widget-info":           {fmt.Sprintf("YA:OneClick/mobile/%s", t.taskPrimeData.ParentWidgetInfo)},
		"origin":                    {"https://www.amazon.com"},
		"user-agent":                {t.taskData.UserAgent},
		// "referer":         {referer},
		"content-type":   {"application/x-www-form-urlencoded; charset=UTF-8"},
		"sec-fetch-dest": {"empty"},
		http.HeaderOrderKey: {
			"host",
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
			"content-type",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqRegisterWidget() (*http.Response, error) {
	path := "https://apx-security.amazon.com/cpe/pm/register"

	data := []KeyValue{
		{Key: "widgetState", Value: t.taskPrimeData.WidgetState},
		{Key: "returnUrl", Value: "/https://www.amazon.com:443/cpe/yourpayments/settings/manageoneclick"},
		{Key: "clientId", Value: "YA:OneClick"},
		{Key: "usePopover", Value: "false"},
		{Key: "maxAgeSeconds", Value: "900"},
		{Key: "iFrameName", Value: "ApxSecureIframe-pp-ZnCW4d-12"},
		{Key: "parentWidgetInstanceId", Value: t.taskPrimeData.ParentWidgetInfo},
		{Key: "hideAddPaymentInstrumentHeader", Value: "true"},
		{Key: "creatablePaymentMethods", Value: "CC"},
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(encodeData(data)))
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host":   {"apx-security.amazon.com"},
		"Accept": {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		// "X-MASH-CSM-TC": {fmt.Sprintf("%d", time.Now().UnixMilli())}, // timestamp e.g. 1711320773636
		"Sec-Fetch-Site":  {"same-site"},
		"Accept-Language": {"en-US,en;q=0.9"},
		"Sec-Fetch-Mode":  {"navigate"},
		"Origin":          {"https://www.amazon.com"},
		"User-Agent":      {t.taskData.UserAgent},
		"Sec-Fetch-Dest":  {"iframe"},
		"Content-Type":    {"application/x-www-form-urlencoded"},
		http.HeaderOrderKey: {
			"host",
			"accept",
			// "x-mash-csm-tc",
			"sec-fetch-site",
			"accept-language",
			"sec-fetch-mode",
			"origin",
			"user-agent",
			"referer",
			"sec-fetch-dest",
			"content-type",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *AmazonTask) tryReqAddPurchasePref() (*http.Response, error) {
	path := "https://www.amazon.com/payments-portal/data/widgets2/v1/customer/" + t.taskPrimeData.CustomerId + "/continueWidget"

	data := []KeyValue{
		{Key: "ppw-widgetEvent:AddOneClickEvent:{}", Value: "+  Add a purchase preference"},
		{Key: "ppw-jsEnabled", Value: "true"},
		{Key: "ppw-widgetState", Value: t.taskPrimeData.WidgetState},
		{Key: "ie", Value: "UTF-8"},
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
		"apx-widget-info":           {fmt.Sprintf("YA:OneClick/mobile/%s", t.taskPrimeData.ParentWidgetInfo)},
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

func (t *AmazonTask) tryReq1ClickPage() (*http.Response, error) {
	path := "https://www.amazon.com/cpe/manageoneclick?_encoding=UTF8&ref_=aw_ya_hp_oneclick_aui"

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Host": {"www.amazon.com"},
		// "x-smashintercepted": {"YES"},
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"sec-fetch-site":  {"cross-site"},
		"accept-language": {"en-US,en;q=0.9"},
		"sec-fetch-mode":  {"navigate"},
		"user-agent":      {t.taskData.UserAgent},
		"referer":         {"https://www.amazon.com/gp/aw/ya?ref_=nav_youraccount_btn&from=hz&isInternal=false"},
		"x-mash-csm-tc":   {fmt.Sprintf("%d", time.Now().UnixMilli())},
		"sec-fetch-dest":  {"document"},
		http.HeaderOrderKey: {
			"host",
			// "x-smashintercepted",
			"accept",
			"sec-fetch-site",
			"accept-language",
			"sec-fetch-mode",
			"user-agent",
			"referer",
			"x-mash-csm-tc",
			"sec-fetch-dest",
		},
	}

	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
