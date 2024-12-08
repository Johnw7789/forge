package request

import (
	"errors"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
)

func (t *AmazonTask) Do2FAStep(step string) error {
	switch step {
	case "init2FA":
		return t.init2FA()
	case "deletePhone":
		return t.deletePhone()
	case "submit2FA":
		return t.submit2FA()
	default:
		return errors.New("Invalid step")
	}
}

func (t *AmazonTask) init2FA() error {
	t.UpdateStatus("Initializing Account Management")

	if t.client == nil {
		err := t.initClient()
		if err != nil {
			return err
		}
	}

	t.UpdateStatus("Requesting Account Management")

	tb := time.Now()

	// * Get the account management page
	_, err := t.reqAccPage()
	if err != nil {
		return err
	}

	t.UpdateStatus("Requesting Security Endpoint")

	// * Get the security page
	bodyStr, resp, err := t.reqSecurityPage()
	if err != nil {
		return err
	}

	// * Check if the page has a verification, if it doesn't then exit early and keep going
	if !t.hasVerification(resp) {
		return nil
	}

	t.UpdateStatus("Getting Account Access OTP")

	code, err := t.ImapClient.FetchOtp(t.UserInfo.Email, tb)
	if err != nil {
		return err
	}

	if code == "" {
		return errors.New("Failed to Fetch OTP for 2FA")
	}

	t.UpdateStatus("Submitting Account Access OTP")

	// * Parse required data for OTP submission
	csrf, err := t.parseSecurityCsrf(bodyStr)
	if err != nil {
		return err
	}

	arb, err := t.parseArbFromParam(resp)
	if err != nil {
		return err
	}

	returnTo, err := t.parseReturnTo(bodyStr)
	if err != nil {
		return err
	}

	// * Submit the OTP code
	_, err = t.reqApproveOTP(code, arb, returnTo, csrf)
	return err
}

func (t *AmazonTask) submit2FA() error {
	t.UpdateStatus("Submitting 2FA")

	// * Request the 2FA page
	bodyStr, err := t.req2FAPage()
	if err != nil {
		return err
	}

	// * Parse the 2FA secret from the response html
	secret, err := t.parse2FASecret(bodyStr)
	if err != nil {
		return err
	}

	twoFactorPayload, err := t.parse2FAPage(bodyStr)
	if err != nil {
		return err
	}

	// * Generate the time-based OTP based on the secret
	otp, err := t.getHOTPToken(secret)
	if err != nil {
		return errors.New("Failed to Generate 2FA OTP")
	}

	_, err = t.req2FASubmit(twoFactorPayload.CsrfToken, twoFactorPayload.SharedSecretId, otp)
	if err != nil {
		return err
	}

	// * Request the 2FA popup
	bodyStr, err = t.req2FAPopup()
	if err != nil {
		return err
	}

	csrfToken, err := t.parseCsrf(bodyStr)
	if err != nil {
		return err
	}

	// * Attempt to finalize the 2FA process
	bodyStr, err = t.req2FAConfirm(csrfToken)
	if err != nil {
		return err
	}

	if !t.completed2FA(bodyStr) {
		return errors.New("Failed to Complete 2FA")
	}

	t.UpdateStatus("2FA Completed")
	t.Secret = secret
	return nil
}

func (t *AmazonTask) deletePhone() error {
	t.UpdateStatus("Initializing Phone Removal")

	// * Delete the phone number from the account if it exists
	bodyStr, err := t.reqPhonePage()
	if err != nil {
		return err
	}

	phoneData, err := t.parsePhonePayload(bodyStr)
	if err != nil {
		return err
	}

	t.UpdateStatus("Requesting Phone Removal")

	bodyStr, err = t.reqDeletePhone(phoneData.AppActionToken, phoneData.AppAction, phoneData.PrevRID, phoneData.WorkflowState)
	if err != nil {
		return err
	}

	if !t.phoneDeleted(bodyStr) {
		t.UpdateStatus("Failed to Remove Phone")
		// * Don't return an error here, the phone number may not exist
		// return errors.New("Failed to Delete Phone")
	} else {
		t.UpdateStatus("Phone Removed")
	}

	return nil
}

func (t *AmazonTask) completed2FA(body string) bool {
	return strings.Contains(body, "You've turned on Two-Step Verification")
}

func (t *AmazonTask) parse2FAPage(body string) (*TwoFactorPayload, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, errors.New("Failed to Parse 2FA (1)")
	}

	csrf, exists := doc.Find("input[name='csrfToken']").Attr("value")
	if !exists {
		return nil, errors.New("Failed to Parse 2FA (2)")
	}

	sharedSecretId, exists := doc.Find("input[name='sharedSecretId']").Attr("value")
	if !exists {
		return nil, errors.New("Failed to Parse 2FA (3)")
	}

	return &TwoFactorPayload{
		CsrfToken:      csrf,
		SharedSecretId: sharedSecretId,
	}, nil
}

func (t *AmazonTask) parse2FASecret(body string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", errors.New("Failed to Parse Secret (1)")
	}

	// <span id="sia-auth-app-formatted-secret">MQOL R5DR OQIX E4RK 6LHL K2JS 7MOJ WADU AJVO RIXL 75PF QI4E X2EA</span>
	secret := doc.Find("#sia-auth-app-formatted-secret").First().Text()
	if secret == "" {
		return "", errors.New("Failed to Parse Secret (2)")
	}

	return secret, nil
}

func (t *AmazonTask) parsePhonePayload(body string) (*PhoneRemovalPayload, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, errors.New("Failed to Parse Phone Payload (1)")
	}

	appActionToken, exists := doc.Find("input[name='appActionToken']").Attr("value")
	if !exists {
		return nil, errors.New("Failed to Parse Phone Payload (2)")
	}

	appAction, exists := doc.Find("input[name='appAction']").Attr("value")
	if !exists {
		return nil, errors.New("Failed to Parse Phone Payload (3)")
	}

	prevRID, exists := doc.Find("input[name='prevRID']").Attr("value")
	if !exists {
		return nil, errors.New("Failed to Parse Phone Payload (4)")
	}

	workflowState, exists := doc.Find("input[name='workflowState']").Attr("value")
	if !exists {
		return nil, errors.New("Failed to Parse Phone Payload (5)")
	}

	return &PhoneRemovalPayload{
		AppActionToken: appActionToken,
		AppAction:      appAction,
		PrevRID:        prevRID,
		WorkflowState:  workflowState,
	}, nil
}

func (t *AmazonTask) parseCsrf(body string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", errors.New("Failed to Parse CSRF (1)")
	}

	csrf, exists := doc.Find("input[name='csrfToken']").Attr("value")
	if !exists {
		return "", errors.New("Failed to Parse CSRF (2)")
	}

	return csrf, nil
}

func (t *AmazonTask) parseReturnTo(body string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", errors.New("Failed to Parse RT (1)")
	}

	returnto, exists := doc.Find("input[name='openid.return_to']").Attr("value")
	if !exists {
		return "", errors.New("Failed to Parse RT (2)")
	}

	return html.UnescapeString(returnto), nil
}

func (t *AmazonTask) parseSecurityCsrf(bodyStr string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		return "", errors.New("Failed to Parse CSRF (1s)")
	}

	csrf, exists := doc.Find("input[name='csrfToken']").Attr("value")
	if !exists {
		return "", errors.New("Failed to Parse CSRF (2s)")
	}

	return csrf, nil
}

func (t *AmazonTask) parseArbFromParam(resp *http.Response) (string, error) {
	location := resp.Request.URL.String()
	if !strings.Contains(location, "arb=") {
		return "", errors.New("Failed to Parse ARB (1)")
	}

	u, err := url.Parse(location)
	if err != nil {
		return "", errors.New("Failed to Parse ARB (2)")
	}

	arb := u.Query().Get("arb")
	if arb == "" {
		return "", errors.New("Failed to Parse ARB (3)")
	}

	return arb, nil
}
