package request

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"html"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Johnw7789/forge/backend/discord"
	"github.com/Johnw7789/forge/backend/shr"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"
)

func (t *AmazonTask) hasVerification(resp *http.Response) bool {
	return strings.Contains(resp.Request.URL.String(), "transactionapproval")
}

func (t *AmazonTask) phoneDeleted(body string) bool {
	return strings.Contains(body, "Mobile number deleted.")
}

func (t *AmazonTask) getHOTPToken(secret string) (string, error) {
	secret = strings.ReplaceAll(secret, " ", "") + "===="

	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", err
	}

	intv := time.Now().Unix() / 30

	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(intv))

	hash := hmac.New(sha1.New, key)
	hash.Write(bs)
	h := hash.Sum(nil)
	o := (h[19] & 15)

	var header uint32
	r := bytes.NewReader(h[o : o+4])
	err = binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return "", err
	}

	h12 := (int(header) & 0x7fffffff) % 1000000
	otp := strconv.Itoa(int(h12))

	return prefixZ(otp), nil
}

func prefixZ(otp string) string {
	if len(otp) == 6 {
		return otp
	}

	for i := (6 - len(otp)); i > 0; i-- {
		otp = "0" + otp
	}

	return otp
}

func encodeData(data []KeyValue) string {
	var encodedData []string
	for _, kv := range data {
		encodedData = append(encodedData, url.QueryEscape(kv.Key)+"="+url.QueryEscape(kv.Value))
	}

	return strings.Join(encodedData, "&")
}

func joinData(data []KeyValue) string {
	var joinedData []string
	for _, kv := range data {
		joinedData = append(joinedData, kv.Key+"="+kv.Value)
	}

	return strings.Join(joinedData, "&")
}

// * Alert the user that the account was created successfully
func (t *AmazonTask) pingDiscord() {
	discord.AlertAccountSuccess(t.DiscordInfo.WebhookSuccess, t.formatProxy(t.UserInfo.Proxy), t.UserInfo.FirstName+" "+t.UserInfo.LastName, t.UserInfo.Email, t.UserInfo.Password)
}

func (t *AmazonTask) parseRegisterUrl(bodyStr string) (string, error) {
	// <a id="createAccountSubmit" href="https://www.amazon.com/ap/register?clientContext=ddf5326dd917382316b434dbf5228c&amp;showRememberMe=true&amp;openid.pape.max_auth_age=43200&amp;openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&amp;pageId=amzn_mturk_requester_v2&amp;openid.return_to=https%3A%2F%2Frequester.mturk.com%2F&amp;prevRID=QVYM5D4D4SJ39DEFCXQ8&amp;openid.assoc_handle=amzn_mturk_requester_v2&amp;openid.mode=checkid_setup&amp;openid.ns.pape=http%3A%2F%2Fspecs.openid.net%2Fextensions%2Fpape%2F1.0&amp;prepopulatedLoginId=&amp;failedSignInCount=0&amp;openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&amp;openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0" class="a-button-text">Create your Amazon account

	// <a id="register_accordion_header" data-csa-c-func-deps="aui-da-a-accordion" data-csa-c-type="widget" data-csa-interaction-events="click" data-action="a-accordion" class="a-accordion-row a-declarative" href="https://www.amazon.com/ap/register?openid.pape.max_auth_age=0&amp;openid.identity=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&amp;accountStatusPolicy=P1&amp;language=en_US&amp;pageId=amzn_dp_project_dee_ios&amp;openid.return_to=https%3A%2F%2Fwww.amazon.com%2Fap%2Fmaplanding&amp;prevRID=ZK8PSQ3KNEMYKEKCH899&amp;openid.assoc_handle=amzn_dp_project_dee_ios&amp;openid.oa2.response_type=code&amp;openid.mode=checkid_setup&amp;openid.ns.pape=http%3A%2F%2Fspecs.openid.net%2Fextensions%2Fpape%2F1.0&amp;openid.ns.oa2=http%3A%2F%2Fwww.amazon.com%2Fap%2Fext%2Foauth%2F2&amp;openid.oa2.code_challenge_method=S256&amp;prepopulatedLoginId=&amp;failedSignInCount=0&amp;openid.oa2.code_challenge=i2l6zQ9LbeJ7fenXvYEqXf3e4wOtEkRW8MP9U&amp;openid.oa2.scope=device_auth_access&amp;openid.claimed_id=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0%2Fidentifier_select&amp;openid.oa2.client_id=device%3A324e4e375446585759314759455244424931584350513539544423413249564c5635564d32573831&amp;openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0" aria-label="" data-csa-c-id="h83id4-dawj1x-oegwi3-i8b1ir"><i class="a-icon a-accordion-radio a-icon-radio-inactive"></i><h5>
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyStr))
	if err != nil {
		return "", errors.New("Failed to Get Register URL (1)")
	}

	// urlString := doc.Find(`a[id="createAccountSubmit"]`).AttrOr("href", "")
	urlString := doc.Find(`a[id="register_accordion_header"]`).AttrOr("href", "")

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", errors.New("Failed to Get Register URL (2)")
	}

	// Unescape HTML escape sequences
	unescapedURL := html.UnescapeString(parsedURL.String())

	return unescapedURL, nil
}

func generateChallengeData(serial string) (ChallengeData, error) {
	clientId := hex.EncodeToString([]byte(serial + "#A2IVLV5VM2W81")) //     amzn alexa
	// clientId := hex.EncodeToString([]byte(serial + "#A3NWHXTQ4EBCZS"))    amzn mshop
	// clientId := hex.EncodeToString([]byte(serial + "#A2825NDLA7WDZV"))    amzn music

	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return ChallengeData{}, err
	}

	verifier := strings.TrimSuffix(base64.URLEncoding.EncodeToString(verifierBytes), "=")

	h := sha256.Sum256([]byte(verifier))

	return ChallengeData{
		ClientId:         clientId,
		Verifier:         verifier,
		VerifierChecksum: strings.TrimSuffix(base64.URLEncoding.EncodeToString(h[:]), "="),
	}, nil
}

func (t *AmazonTask) accountCreated(name string, body string) (bool, error) {
	if strings.Contains(body, name) {
		return true, nil
	}

	return false, nil
}

func (t *AmazonTask) evalVerifyEmail(resp http.Response) (bool, error) {
	if strings.Contains(resp.Request.URL.String(), "new_account=1") {
		return true, nil
	}

	return false, nil
}

func (t *AmazonTask) evalCreateAccount(body string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return "", errors.New("Failed to Get VTI")
	}

	return doc.Find(`input[name="verifyToken"]`).AttrOr("value", ""), nil
}

func (t *AmazonTask) evalCreateData(body string) (CreateInfo, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return CreateInfo{}, err
	}

	// ap_register_form
	form := doc.Find(`form[id="ap_register_form"]`)

	ld := CreateInfo{
		AppAction:      form.Find(`input[name="appAction"]`).AttrOr("value", ""),
		AppActionToken: form.Find(`input[name="appActionToken"]`).AttrOr("value", ""),
		PrevRID:        form.Find(`input[name="prevRID"]`).AttrOr("value", ""),
		WorkflowState:  form.Find(`input[name="workflowState"]`).AttrOr("value", ""),
		OpenIdReturnTo: form.Find(`input[name="openid.return_to"]`).AttrOr("value", ""),
	}

	return ld, nil
}

func (t *AmazonTask) getInnerHTML(reader *bytes.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return "", err
	}

	// Find the body tag and get its inner HTML
	return doc.Find("body").Html()
}

func (t *AmazonTask) formatProxy(proxy shr.Proxy) string {
	return strings.Join([]string{proxy.Host, proxy.Port, proxy.User, proxy.Pass}, ":")
}

func (t *AmazonTask) getBodyStr(resp *http.Response) (string, error) {
	if resp == nil {
		return "", errors.New("Invalid Response!")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (t *AmazonTask) httpPostHeaders() http.Header {
	return http.Header{
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
			"host",
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
}

func (t *AmazonTask) httpGetHeaders() http.Header {
	return http.Header{
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"sec-fetch-site":  {"none"},
		"sec-fetch-mode":  {"navigate"},
		"user-agent":      {t.taskData.UserAgent},
		"accept-language": {"en-US,en;q=0.9"},
		"sec-fetch-dest":  {"document"},
		http.HeaderOrderKey: {
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"user-agent",
			"accept-language",
			"sec-fetch-dest",
		},
	}
}

func (t *AmazonTask) httpGetSafariHeaders() http.Header {
	return http.Header{
		"Host":            {"www.amazon.com"},
		"accept":          {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"sec-fetch-site":  {"none"},
		"sec-fetch-mode":  {"navigate"},
		"user-agent":      {t.taskData.UserAgent},
		"accept-language": {"en-US,en;q=0.9"},
		"sec-fetch-dest":  {"document"},
		http.HeaderOrderKey: {
			"Host",
			"Cookie",
			"accept",
			"sec-fetch-site",
			"sec-fetch-mode",
			"user-agent",
			"accept-language",
			"sec-fetch-dest",
		},
	}
}
