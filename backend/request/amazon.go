package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/Johnw7789/forge/backend/antibot"
	"github.com/Johnw7789/forge/backend/sms/daisysms"
	"github.com/Johnw7789/forge/backend/sms/smsman"
	"github.com/Johnw7789/forge/backend/sms/smspool"

	b64 "encoding/base64"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

// * Cycle through the steps of the task for as long as it is running on the frontend side
func (t *AmazonTask) DoStep(step string) error {
	switch step {
	case "init":
		return t.init()
	case "getCreateData":
		return t.getCreateData()
	case "submitCreate":
		return t.submitCreate()
	case "pingDiscord":
		if t.taskData.Success {
			t.pingDiscord()
		}
		return nil
	case "phone2FA":
		var steps = []string{"init2FA", "deletePhone", "submit2FA"}
		for _, s := range steps {
			if s == "deletePhone" && t.Phone == "" {
				continue
			}

			err := t.Do2FAStep(s)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New("Invalid step")
	}
}

// * Initialize the task, create the http client and the sms client depending on the user's sms provider
func (t *AmazonTask) init() error {
	t.UpdateStatus("Initializing")
	err := t.initClient()
	if err != nil {
		return err
	}

	switch t.UserInfo.SmsInfo.Provider {
	case "SMS Man":
		t.smClient, err = smsman.NewSMSManClient(t.UserInfo.SmsInfo.ApiKey)
		if err != nil {
			t.UpdateStatus("Failed to Create SMS Client")
		}
	case "SMS Pool":
		t.spClient, err = smspool.NewSMSPoolClient(t.UserInfo.SmsInfo.ApiKey)
		if err != nil {
			t.UpdateStatus("Failed to Create SMS Client")
		}
	case "Daisy SMS":
		t.dsClient, err = daisysms.NewDaisySMSClient(t.UserInfo.SmsInfo.ApiKey)
		if err != nil {
			t.UpdateStatus("Failed to Create SMS Client")
		}
	default:
		t.UpdateStatus("Invalid SMS Provider")
		return errors.New("Invalid SMS Provider")
	}

	t.cc = antibot.NewCaptchaClient(t.UserInfo.CaptchaInfo.APIKey)

	return err
}

func (t *AmazonTask) initClient() error {
	jar := tls_client.NewCookieJar()

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithCookieJar(jar),
		tls_client.WithClientProfile(profiles.Safari_IOS_17_0),
	}

	t.taskData.UserAgent, t.taskData.OSVer = antibot.NewIOS17UserAgent()

	if t.UserInfo.Proxy.Host != "" && t.UserInfo.Proxy.Port != "" {
		url := ""

		if t.UserInfo.Proxy.User != "" && t.UserInfo.Proxy.Pass != "" {
			url = "http://" + t.UserInfo.Proxy.User + ":" + t.UserInfo.Proxy.Pass + "@" + t.UserInfo.Proxy.Host + ":" + t.UserInfo.Proxy.Port
		} else {
			url = "http://" + t.UserInfo.Proxy.Host + ":" + t.UserInfo.Proxy.Port
		}

		options = append(options, tls_client.WithProxyUrl(url))
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		return err
	}

	t.client = client

	return nil
}

func (t *AmazonTask) getCreateData() error {
	t.UpdateStatus("Initializing")

	// * MapMD cookie, differs from each app, we are using the Alexa app hence the bundleid
	mapMdCookie := MapMD{
		DeviceUserDictionary: []interface{}{},
		DeviceRegistrationData: DeviceRegistrationData{
			SoftwareVersion: "1",
		},
		AppIdentifier: AppIdentifier{
			AppVersion: "2.2.595606",
			BundleID:   "com.amazon.echo",
		},
	}

	device := antibot.GetRandAlexaDevice(t.taskData.OSVer)
	secureCookie, err := antibot.GenerateSecureCookie(device)
	if err != nil {
		t.UpdateStatus("Failed to Generate Secure Cookie")
		return err
	}

	// * Marshal, and encode the cookie to base64
	jsonMarshalled, err := json.Marshal(mapMdCookie)
	if err != nil {
		return errors.New("Failed to Marshal MapMD")
	}

	b64MapMD := b64.StdEncoding.EncodeToString(jsonMarshalled)

	cookies := []*http.Cookie{
		{
			Name:  "map-md",
			Value: b64MapMD,
		},
		{
			Name:  "frc",
			Value: secureCookie,
		},
	}

	u, err := url.Parse("https://www.amazon.com")

	if err == nil {
		t.client.SetCookies(u, cookies)
	}

	cd, err := generateChallengeData(device.Serial)
	if err != nil {
		return err
	}

	t.taskData.ChallengeData = cd

	// * Init the session by requesting the signin page
	bodyStr, signinUrl, err := t.reqSigninAlexa()
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * time.Duration(rand.Intn(3000)+1000))

	// * Parse the register (signup) url from the signin page
	registerUrl, err := t.parseRegisterUrl(bodyStr)
	if err != nil {
		return err
	}

	// * Get the creation data required for the account creation
	bodyStr, err = t.reqCreateDataAlexa(registerUrl, signinUrl)
	if err != nil {
		return err
	}

	t.taskData.CreateInfo, err = t.evalCreateData(bodyStr)
	if err != nil {
		return err
	}

	sleepMs := rand.Intn(20000) + 10000

	status := fmt.Sprintf("Sleeping for %d Seconds", sleepMs/1000)
	t.UpdateStatus(status)

	time.Sleep(time.Millisecond * time.Duration(sleepMs))

	return nil
}

func (t *AmazonTask) submitCreate() error {
	// * Mark the time so that we can get a OTP code that is AFTER this time and not an old one
	tb := time.Now()

	t.UpdateStatus("Creating Account")

	verifyToken := ""
	// * Request the account submit endpoint
	bodyStr, err := t.reqCreateAccountAlexa()
	if err != nil {
		return err
	}

	// * Check if the account creation has a captcha
	hasCaptcha := t.hasCaptcha(bodyStr)
	if hasCaptcha {
		t.UpdateStatus("Bypassing Puzzle")

		for i := 0; i < t.UserInfo.CaptchaInfo.MaxRetries; i++ {
			bodyStr, err = t.HandleCaptcha(bodyStr)
			if err != nil {
				continue
			}

			if !t.hasCaptcha(bodyStr) {
				break
			}
		}

		// * If the captcha is still present after the max retries, return an error
		if t.hasCaptcha(bodyStr) {
			return errors.New("Failed to Bypass Puzzle")
		}
		// * Check if the account creation has a puzzle, which is a different much harder type of captcha that is not solvable through third party services
	} else if t.hasPuzzle(bodyStr) {
		return errors.New("Invalid Flow - Error")
	}

	// * Something went wrong on the Amazon side, return an error
	if strings.Contains(bodyStr, "Internal Error. Please try again later.") {
		return errors.New("Amazon: Internal Error")
	}

	t.UpdateStatus("Fetching Email OTP")

	// * Attempt to fetch the OTP code from the email
	code, err := t.ImapClient.FetchOtp(t.UserInfo.Email, tb)
	if err != nil {
		return err
	}

	if code == "" {
		return errors.New("Failed to Fetch Email OTP")
	}

	t.UpdateStatus("Found Email OTP: " + code)

	verifyToken, err = t.evalCreateAccount(bodyStr)
	if err != nil {
		return err
	}

	// * Submit the OTP code to verify the email
	bodyStr, err = t.reqVerifyEmailMobile(code, verifyToken)
	if err != nil {
		return err
	}

	newToken, err := t.evalCreateAccount(bodyStr)
	if err != nil {
		return err
	}

	if newToken != "" {
		verifyToken = newToken
	}

	// * Check if the account creation requires a phone number
	if strings.Contains(bodyStr, "Add mobile number") || hasCaptcha {
		err = t.submitSMS(t.UserInfo.SmsInfo.Provider, verifyToken)
		if err != nil {
			return err
		}
	}

	bodyStr, err = t.reqHomePage()
	if err != nil {
		return err
	}

	// * Check if the account was created successfully
	accountCreated, err := t.accountCreated(t.UserInfo.FirstName, bodyStr)
	if err != nil {
		return err
	} else if !accountCreated {
		return errors.New("Unable to Verify Account Creation")
	}

	// * Get the cookies from the client
	var cookies string
	u, _ := url.Parse("https://www.amazon.com")
	for _, cookie := range t.client.GetCookies(u) {
		cookies = cookies + cookie.Name + "=" + cookie.Value + ";"
	}
	cookies = strings.TrimSuffix(cookies, ";")
	t.Cookies = cookies

	t.taskData.Success = true
	t.UpdateStatus("Account Created")

	return nil
}
