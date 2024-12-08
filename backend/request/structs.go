package request

import (
	"github.com/Johnw7789/forge/backend/antibot"
	"github.com/Johnw7789/forge/backend/email"
	"github.com/Johnw7789/forge/backend/shr"
	"github.com/Johnw7789/forge/backend/sms/daisysms"
	"github.com/Johnw7789/forge/backend/sms/smsman"
	"github.com/Johnw7789/forge/backend/sms/smspool"

	tls_client "github.com/bogdanfinn/tls-client"
)

type AmazonTask struct {
	UserInfo    UserInfo
	DiscordInfo DiscordInfo
	FP          shr.Fingerprint

	Payment Payment
	Address Address

	taskData      TaskData
	taskPrimeData TaskPrimeData

	client     tls_client.HttpClient
	ImapClient *email.EmailClient

	spClient smspool.SMSPoolClient
	smClient smsman.SMSManClient
	dsClient daisysms.DaisySMSClient
	cc       antibot.CaptchaClient

	Phone     string
	Secret    string
	ChromeVer string
	Cookies   string

	UpdateStatus func(string)
	BadFunc      tls_client.BadPinHandlerFunc
}

type TaskData struct {
	Metadata      string
	UserAgent     string
	OSVer         string
	Referer       string
	Location      string
	ChallengeData ChallengeData
	CreateInfo    CreateInfo
	Success       bool
}

// "YA:OneClick/mobile/%s"
type TaskPrimeData struct {
	WidgetInfo       string
	ParentWidgetInfo string
	WidgetState      string
	FinalWidgetState string
	InstrumentId     string
	CustomerId       string
	SessionId        string
	PaymentMethodId  string
	AddressId        string
}

type PrimeSignupData struct {
	CampaignId                 string
	OfferId                    string
	LocationID                 string
	OfferToken                 string
	RedirectURL                string
	CancelRedirectURL          string
	PreviousContainerRequestId string
	ActionPageDefinitionId     string
}

type CaptchaInfo struct {
	CaptchaUrl    string
	CaptchaToken  string
	CaptchaType   string
	ClientContext string
	CodeChallenge string
	ClientId      string
	VerifyToken   string
	CaptchaInput  string
}

type UserInfo struct {
	Email       string
	Password    string
	FirstName   string
	LastName    string
	Proxy       shr.Proxy
	ImapInfo    shr.ImapInfo
	SmsInfo     shr.SMSInfo
	CaptchaInfo shr.CaptchaInfo
}

type Payment struct {
	CardHolder string
	CardNumber string
	ExpMonth   string
	ExpYear    string
	CVV        string
}

type Address struct {
	FullName    string
	Line1       string
	Line2       string
	City        string
	State       string
	Zip         string
	PhoneNumber string
}

type DiscordInfo struct {
	User           string
	WebhookSuccess string
	WebhookFail    string
}

type CreateInfo struct {
	AppAction      string
	AppActionToken string
	PrevRID        string
	WorkflowState  string
	OpenIdReturnTo string
}

type ChallengeData struct {
	ClientId         string
	Verifier         string
	VerifierChecksum string
}

type MapMD struct {
	DeviceUserDictionary   []interface{}          `json:"device_user_dictionary"`
	DeviceRegistrationData DeviceRegistrationData `json:"device_registration_data"`
	AppIdentifier          AppIdentifier          `json:"app_identifier"`
}

type DeviceRegistrationData struct {
	SoftwareVersion string `json:"software_version"`
}

type AppIdentifier struct {
	AppVersion string `json:"app_version"`
	BundleID   string `json:"bundle_id"`
}

type TwoFactorPayload struct {
	CsrfToken      string
	SharedSecretId string
}

type PhoneRemovalPayload struct {
	AppActionToken string
	AppAction      string
	PrevRID        string
	WorkflowState  string
}

type KeyValue struct {
	Key   string
	Value string
}
