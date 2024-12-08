package request

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *AmazonTask) HandleCaptcha(bodyStr string) (string, error) {
	captchaInfo, err := t.getCaptchaInfo(bodyStr)
	if err != nil {
		return "", err
	}

	captchaInfo.ClientId = t.taskData.ChallengeData.ClientId
	captchaInfo.CodeChallenge = t.taskData.ChallengeData.Verifier

	img, err := t.reqCaptchaImage(captchaInfo.CaptchaUrl)
	if err != nil {
		return "", err
	}

	code, err := t.cc.SolveCap(img)
	if err != nil {
		return "", err
	}

	captchaInfo.CaptchaInput = code

	return t.reqCaptchaSubmit(captchaInfo)
}

func (t *AmazonTask) hasPuzzle(respBody string) bool {
	return strings.Contains(respBody, "setupACIC")
}

func (t *AmazonTask) hasCaptcha(respBody string) bool {
	return strings.Contains(respBody, "Enter the letters and numbers above")
}

func (t *AmazonTask) getCaptchaInfo(body string) (CaptchaInfo, error) {
	var ci CaptchaInfo

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return ci, err
	}

	ci.CaptchaUrl, _ = doc.Find("img[alt=captcha]").Attr("src")
	ci.CaptchaToken, _ = doc.Find("input[name=cvf_captcha_captcha_token]").Attr("value")
	ci.CaptchaType, _ = doc.Find("input[name=cvf_captcha_captcha_type]").Attr("value")
	ci.ClientContext, _ = doc.Find("input[name=clientContext]").Attr("value")
	ci.VerifyToken, _ = doc.Find("input[name=verifyToken]").Attr("value")

	return ci, nil
}
