package shr

import (
	"sync"
)

type Settings struct {
	LicenseKey    string   `json:"licenseKey"`
	ProxiesFile   string   `json:"proxiesFile"`
	EmailFile     string   `json:"emailFile"`
	Webhook       string   `json:"webhook"`
	ImapInfo      ImapInfo `json:"imapInfo"`
	SMSInfo       *SMSInfo `json:"smsInfo"`
	TwoCaptchaKey string   `json:"twoCaptchaKey"`
	MaxTasks      int      `json:"maxTasks"`
	LocalHost     bool     `json:"localhost"`
}

type AccountOutput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Proxy     string `json:"proxy"`
}

type PaymentData struct {
	CardName string
	CardNum  string
	CardExpM string
	CardExpD string
}

type ShippingData struct {
	Address1 string
	Address2 string
	ZipCode  string
	City     string
	State    string
}

type Proxy struct {
	Host string
	Port string
	User string
	Pass string
}

type ImapInfo struct {
	ImapMu   *sync.Mutex
	Username string `json:"username"`
	Password string `json:"password"`
}

type SMSInfo struct {
	SMSMu      *sync.Mutex
	Username   string
	ApiKey     string
	Provider   string `json:"provider"`
	MaxRetries int    `json:"maxRetries"`
}

type CaptchaInfo struct {
	APIKey     string
	MaxRetries int
}

type AccountResult struct {
	Success       bool
	Email         string
	Password      string
	JiggedAddress string
}

type Canvas struct {
	FP            string `json:"fp"`
	EmailFP       string `json:"emailFp"`
	HistogramBins []int  `json:"histogramBins"`
}

type Math struct {
	Sin string `json:"sin"`
	Cos string `json:"cos"`
	Tan string `json:"tan"`
}

type Screen struct {
	Height      string
	Width       string
	AvailHeight string
	Colordepth  string
}

type Fingerprint struct {
	GPUVendor     string   `json:"gpuVendor"`
	GPURenderer   string   `json:"gpuRenderer"`
	GPUExtensions []string `json:"gpuExtensions"`
	Screen        Screen
	Canvas        Canvas
	Math          Math `json:"math"`
}
