package request

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
	"github.com/tidwall/gjson"
	"golang.design/x/clipboard"
)

func (t *AmazonTask) DoInfoStep(step string) error {
	switch step {
	case "init":
		return t.initInfo()
	case "submitAddress":
		return t.submitAddress()
		// return t.forceAddress()
	case "submitPayment":
		return t.submitPayment()
	case "submitProfile":
		return t.submitProfile()
	default:
		return errors.New("Invalid step")
	}
}

func (t *AmazonTask) DoPrimeStep(step string) error {
	switch step {
	case "initPrime":
		return t.initPrime()
	case "submitPrime":
		return t.submitPrime()
	default:
		return errors.New("Invalid step")
	}
}

func (t *AmazonTask) initInfo() error {
	t.UpdateStatus("Initializing Info Submit")
	if t.client == nil {
		err := t.initClient()
		if err != nil {
			return err
		}

		u, _ := url.Parse("https://www.amazon.com")
		u2, _ := url.Parse("https://apx-security.amazon.com")

		cookiesSpl := strings.Split(t.Cookies, ";")
		var cookies []*http.Cookie
		for _, cookieStr := range cookiesSpl {
			cookieSpl := strings.Split(cookieStr, "=")
			if len(cookieSpl) != 2 {
				continue
			}

			cookies = append(cookies, &http.Cookie{
				Name:  cookieSpl[0],
				Value: cookieSpl[1],
			})
		}

		t.client.SetCookies(u, cookies)
		t.client.SetCookies(u2, cookies)
	}

	// Get the initial page
	bodyStr, err := t.req1ClickPage()
	if err != nil {
		return err
	}

	var instanceId string
	instanceId, t.taskPrimeData.WidgetState, _, t.taskPrimeData.CustomerId, err = t.parseOptions(bodyStr)
	if err != nil {
		return err
	}

	// t.taskPrimeData.ParentWidgetInfo = fmt.Sprintf("YA:OneClick/mobile/%s", instanceId)
	t.taskPrimeData.ParentWidgetInfo = instanceId

	return err
}

func (t *AmazonTask) initPrime() error {
	t.UpdateStatus("Initializing Prime Signup")
	if t.client == nil {
		err := t.initClient()
		if err != nil {
			return err
		}

		u, _ := url.Parse("https://www.amazon.com")

		cookiesSpl := strings.Split(t.Cookies, ";")
		var cookies []*http.Cookie
		for _, cookieStr := range cookiesSpl {
			cookieSpl := strings.Split(cookieStr, "=")
			if len(cookieSpl) != 2 {
				continue
			}

			cookies = append(cookies, &http.Cookie{
				Name:  cookieSpl[0],
				Value: cookieSpl[1],
			})
		}

		t.client.SetCookies(u, cookies)
	}

	// Get the initial page
	_, err := t.reqHomePage()
	if err != nil {
		return err
	}

	return err
}

func (t *AmazonTask) submitAddress() error {
	t.UpdateStatus("Submitting Address")

	bodyStr, err := t.reqAddPurchasePref()
	if err != nil {
		return err
	}

	// Parse the widget state from the response
	t.taskPrimeData.WidgetState, err = t.parseWidgetStateHtml(bodyStr)
	if err != nil {
		return err
	}

	bodyStr, err = t.reqAddAddress()
	if err != nil {
		return err
	}

	t.taskPrimeData.WidgetState, err = t.parseWidgetStateHtml(bodyStr)
	if err != nil {
		return err
	}

	if !strings.Contains(bodyStr, "Add a payment method") {
		bodyStr, err = t.reqAddAddress()
		if err != nil {
			return err
		}

		t.taskPrimeData.WidgetState, err = t.parseWidgetStateHtml(bodyStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *AmazonTask) submitPayment() error {
	t.UpdateStatus("Submitting Payment")

	t.taskPrimeData.FinalWidgetState = t.taskPrimeData.WidgetState

	bodyStr, err := t.reqRegisterWidget()
	if err != nil {
		return err
	}

	// Parse the options
	var instanceId string
	instanceId, t.taskPrimeData.WidgetState, _, _, err = t.parseOptions(bodyStr)
	if err != nil {
		return err
	}

	// t.taskPrimeData.WidgetInfo = fmt.Sprintf("YA:OneClick/mobile/%s", instanceId)
	t.taskPrimeData.WidgetInfo = instanceId

	time.Sleep(time.Second * 3)

	bodyStr, err = t.reqCardType()
	if err != nil {
		return err
	}

	cardType := gjson.Get(bodyStr, "additionalWidgetResponseData.additionalData.issuer").String()

	time.Sleep(time.Second * 4)

	bodyStr, err = t.reqSubmitPayment(cardType)
	if err != nil {
		return err
	}

	t.taskPrimeData.WidgetState, err = t.parseWidgetStateHtml(bodyStr)
	if err != nil {
		return err
	}

	return nil
}

func (t *AmazonTask) submitProfile() error {
	t.UpdateStatus("Submitting Address Select")

	bodyStr, err := t.reqAddressId()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 3)

	t.taskPrimeData.AddressId = gjson.Get(bodyStr, "customerIntent.addressId").String()

	bodyStr, err = t.reqAddressSelect()
	if err != nil {
		return err
	}

	success := gjson.Get(bodyStr, "additionalWidgetResponseData.additionalData.widgetDone").Bool()
	if !success {
		return errors.New("Failed to Submit Address Select")
	}

	t.taskPrimeData.InstrumentId = gjson.Get(bodyStr, "additionalWidgetResponseData.additionalData.paymentInstrumentId").String()

	t.UpdateStatus("Submitting Profile")

	_, err = t.reqFinalize1()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 2)

	bodyStr, err = t.reqFinalize2()
	if err != nil {
		return err
	}

	t.taskPrimeData.WidgetState, err = t.parseWidgetStateHtml(bodyStr)
	if err != nil {
		return err
	}

	_, err = t.reqFinalize3()
	if err != nil {
		return err
	}

	t.UpdateStatus("Submitted Profile")

	return nil
}

func (t *AmazonTask) submitPrime() error {
	t.UpdateStatus("Submitting Prime")

	bodyStr, err := t.reqPrimeSignup()
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 5)

	// Parse the prime signup data
	pd, err := t.parsePrimeSignupData(bodyStr)
	if err != nil {
		return err
	}

	bodyStr, err = t.reqPrimeSubmit(pd)
	if err != nil {
		return err
	}

	actionPageDefinitionId, err := t.parsePageId(bodyStr)
	if err != nil {
		return err
	}

	pd.ActionPageDefinitionId = actionPageDefinitionId

	// parse the options
	t.taskPrimeData.WidgetInfo, t.taskPrimeData.WidgetState, t.taskPrimeData.SessionId, t.taskPrimeData.CustomerId, err = t.parseOptions(bodyStr)
	if err != nil {
		return err
	}

	bodyStr, err = t.reqPrimeWidget()
	if err != nil {
		return err
	}

	// bodyStr = strings.ReplaceAll(bodyStr, `\"`, `"`)

	paymentId := gjson.Get(bodyStr, "additionalWidgetResponseData.additionalData.preferencePaymentMethodIds").String()
	paymentId = strings.ReplaceAll(paymentId, `[`, ``)
	paymentId = strings.ReplaceAll(paymentId, `]`, ``)
	paymentId = strings.ReplaceAll(paymentId, `"`, ``)

	if paymentId == "" {
		return errors.New("Failed to Get PID")
	}

	t.taskPrimeData.PaymentMethodId = paymentId

	bodyStr, err = t.reqPrimeFinalize(pd)
	if err != nil {
		return err
	}

	clipboard.Init()
	clipboard.Write(clipboard.FmtText, []byte(bodyStr))

	t.UpdateStatus("Activated Prime")

	return nil
}

func (t *AmazonTask) parseWidgetStateHtml(bodyStr string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(gjson.Get(bodyStr, "htmlContent").String()))
	if err != nil {
		return "", errors.New("Failed to Parse WSH")
	}

	// Find the input element with name ppw-widgetState and get its value
	return doc.Find("input[name='ppw-widgetState']").AttrOr("value", ""), nil
}

func (t *AmazonTask) parseOptions(bodyStr string) (string, string, string, string, error) {
	re := regexp.MustCompile(`var options = ({.*?});`)

	matches := re.FindStringSubmatch(bodyStr)
	if len(matches) < 2 {
		return "", "", "", "", errors.New("Opt Not Found")
	}

	optionsJSON := matches[1]

	widgetInstanceId := gjson.Get(optionsJSON, "widgetInstanceId").String()
	widgetState := gjson.Get(optionsJSON, "serializedState").String()
	customerId := gjson.Get(optionsJSON, "customerId").String()
	sessionId := gjson.Get(optionsJSON, "sessionId").String()

	return widgetInstanceId, widgetState, sessionId, customerId, nil
}





