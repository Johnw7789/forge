package antibot

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/tidwall/gjson"
)

type CaptchaClient struct {
	cclient *api2captcha.Client
}

type AWSCaptchaData struct {
	Iv              string
	SiteKey         string
	Url             string
	Context         string
	ChallengeScript string
	CaptchaScript   string
}

func NewCaptchaClient(apiKey string) CaptchaClient {
	cclient := api2captcha.NewClient(apiKey)

	cclient.PollingInterval = 2

	return CaptchaClient{
		cclient: cclient,
	}
}

func (cc *CaptchaClient) SolveCap(img []byte) (string, error) {
	b64 := base64.StdEncoding.EncodeToString(img)

	cap := api2captcha.Normal{
		Base64: b64,
	}

	code, err := cc.cclient.Solve(cap.ToRequest())
	if err != nil {
		return "", err
	}

	return code, nil
}

func (cc *CaptchaClient) SolveAWSCaptcha(cd AWSCaptchaData) (string, error) {
	cap := api2captcha.AmazonWAF{
		Iv:            cd.Iv,
		SiteKey:       cd.SiteKey,
		Url:           cd.Url,
		Context:       cd.Context,
		CaptchaScript: cd.CaptchaScript,
	}

	tokenCh := make(chan string)
	done := make(chan struct{})
	defer close(done)

	for i := 0; i < 3; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered:", r)
				}
			}()

			resp, err := cc.cclient.Solve(cap.ToRequest())
			if err != nil {
				return
			}

			token := gjson.Get(resp, "captcha_voucher").String()
			if token != "" {
				select {
				case tokenCh <- token:
				case <-done:
					return
				}
			}
		}()
	}

	timeout := time.After(60 * time.Second) // Timeout after 60 seconds, which is when the captcha will expire
	select {
	case token := <-tokenCh:
		return token, nil
	case <-timeout:
		return "", errors.New("timeout: no token found")
	}
}
