package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"os"
	"sync"

	http "github.com/bogdanfinn/fhttp"

	forgeIcloud "github.com/Johnw7789/Go-iClient/icloud"
	forgeDiscord "github.com/Johnw7789/forge/backend/discord"
	forgeImap "github.com/Johnw7789/forge/backend/email"
	"github.com/Johnw7789/forge/backend/request"
	"github.com/Johnw7789/forge/backend/shr"
	forgeShr "github.com/Johnw7789/forge/backend/shr"
	forgeDs "github.com/Johnw7789/forge/backend/sms/daisysms"
	forgeSm "github.com/Johnw7789/forge/backend/sms/smsman"
	forgeSp "github.com/Johnw7789/forge/backend/sms/smspool"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type BackgroundController struct {
	ctx context.Context

	frontendTasks []TaskState

	iclient    *forgeIcloud.Client
	iMu        sync.Mutex
	generating bool
	otpChan    chan string

	imapClient *forgeImap.EmailClient

	tasks     map[string]Task
	accTasks  map[string]Task
	emails    map[string]bool
	proxies   map[string]bool
	accounts  []Account
	cards     []Card
	addresses []Address

	tasksMu     sync.Mutex
	emailsMu    sync.Mutex
	proxiesMu   sync.Mutex
	accountsMu  sync.Mutex
	cardsMu     sync.Mutex
	addressesMu sync.Mutex

	dataMu sync.Mutex

	sc *SettingsController
}

func NewBackgroundController(settingsController *SettingsController) *BackgroundController {
	tasks := make(map[string]Task)
	fTasks := []TaskState{}
	emails := make(map[string]bool)
	proxies := make(map[string]bool)
	accTasks := make(map[string]Task)
	otpChan := make(chan string)

	return &BackgroundController{sc: settingsController, iMu: sync.Mutex{}, tasksMu: sync.Mutex{}, dataMu: sync.Mutex{}, emailsMu: sync.Mutex{}, proxiesMu: sync.Mutex{}, frontendTasks: fTasks, tasks: tasks, emails: emails, proxies: proxies, accountsMu: sync.Mutex{}, cardsMu: sync.Mutex{}, addressesMu: sync.Mutex{}, accTasks: accTasks, otpChan: otpChan}
}

// * Wait for an event to trigger the startup function from the frontend
func (bc *BackgroundController) startup(ctx context.Context) {
	bc.ctx = ctx

	runtime.EventsOn(ctx, "frontend:init", func(a ...interface{}) {
		bc.LoadEmails()
		bc.LoadProxies()
		bc.LoadAccounts()
		bc.LoadCards()
		bc.LoadAddresses()

		if bc.sc.settings.IcloudConfig.Username != "" && bc.sc.settings.IcloudConfig.Password != "" {
			var err error
			bc.iclient, err = forgeIcloud.NewClient(bc.sc.settings.IcloudConfig.Username, bc.sc.settings.IcloudConfig.Password, false)
			if err != nil {
				runtime.EventsEmit(bc.ctx, "error", "Failed to init icloud!")
			}
		}

		if bc.sc.settings.ImapConfig.Username != "" && bc.sc.settings.ImapConfig.Password != "" {
			var err error
			bc.imapClient, err = forgeImap.NewEmailClient(bc.sc.settings.ImapConfig.Username, bc.sc.settings.ImapConfig.Password)
			if err != nil {
				runtime.EventsEmit(bc.ctx, "error", "Failed to init imap!")
			} else {
				runtime.EventsEmit(bc.ctx, "success", "Successfully logged in to imap!")
			}
		}
	})
}

func (bc *BackgroundController) getUnusedEmail() string {
	bc.emailsMu.Lock()
	defer bc.emailsMu.Unlock()

	for email, used := range bc.emails {
		if bc.emailExistsInAccounts(email) {
			continue
		}

		if !used {
			bc.emails[email] = true
			return email
		}
	}

	return ""
}

func (bc *BackgroundController) getProxy() string {
	bc.proxiesMu.Lock()
	defer bc.proxiesMu.Unlock()

	for proxy, used := range bc.proxies {
		if bc.proxyExistsInAccounts(proxy) && bc.sc.settings.LimitProxyUse {
			continue
		}

		if !used {
			bc.proxies[proxy] = true
			return proxy
		}
	}

	return ""
}

func (bc *BackgroundController) getAddressFromProfileId(profileId string) Address {
	bc.addressesMu.Lock()
	defer bc.addressesMu.Unlock()

	for _, address := range bc.addresses {
		if address.Id == profileId {
			return address
		}
	}

	return Address{}
}

func (bc *BackgroundController) getCardFromProfileId(profileId string) Card {
	bc.cardsMu.Lock()
	defer bc.cardsMu.Unlock()

	for _, card := range bc.cards {
		if card.Id == profileId {
			return card
		}
	}

	return Card{}
}

// * Initializes a "add prime" task and immediately starts it
func (bc *BackgroundController) CreatePrimeTask(cookies, proxy, accId string) {
	// fmt.Println("cookies: ", cookies)
	// fmt.Println("proxy: ", proxy)
	// fmt.Println("accId: ", accId)

	// check if task exists with acc id
	bc.accountsMu.Lock()
	if _, exists := bc.accTasks[accId]; exists {
		// emit err event
		runtime.EventsEmit(bc.ctx, "error", "Sub task already in process for account "+accId)
		bc.accountsMu.Unlock()
		return
	}
	bc.accountsMu.Unlock()

	task := Task{
		amzTask: request.AmazonTask{
			UserInfo: request.UserInfo{
				Proxy: forgeShr.Proxy{},
			},
			Cookies: cookies,
		},
	}

	parsedProxy, err := forgeShr.ParseProxyString(proxy)
	if err == nil {
		task.amzTask.UserInfo.Proxy = parsedProxy
	}

	updateStatus := func(status string) {
		// return if task stopped
		if task.stopped == nil {
			return
		}

		// emit status event
		bc.accountsMu.Lock()
		for i, _ := range bc.accounts {
			newAcc := bc.accounts[i]

			if newAcc.Id == accId {
				newAcc.Status = status
				bc.accounts[i] = newAcc
			}
		}
		runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
		bc.accountsMu.Unlock()
	}

	task.amzTask.UpdateStatus = updateStatus

	go func(id string) {
		success, err := task.StartPrime()
		bc.accountsMu.Lock()
		delete(bc.accTasks, id)
		bc.accountsMu.Unlock()
		if err != nil {
			bc.accountsMu.Lock()

			for i, _ := range bc.accounts {
				newAcc := bc.accounts[i]

				if newAcc.Id == id {
					newAcc.Status = "Error: " + err.Error()
					bc.accounts[i] = newAcc
				}
			}

			runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
			bc.accountsMu.Unlock()
		}

		if success {
			bc.accountsMu.Lock()
			for i, _ := range bc.accounts {
				newAcc := bc.accounts[i]

				if newAcc.Id == id {
					newAcc.Status = "Prime Activated"
					newAcc.Prime = true
					bc.accounts[i] = newAcc
				}
			}
			runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
			bc.accountsMu.Unlock()

			bc.SaveAccounts(bc.accounts)
		}
	}(accId)
}

// * Initializes an "add address" task and immediately starts it
func (bc *BackgroundController) CreateInfoTask(addressProfileId, cardProfileId, cookies, proxy, accId string) {
	// check if task exists with acc id
	bc.accountsMu.Lock()
	if _, exists := bc.accTasks[accId]; exists {
		// emit err event
		runtime.EventsEmit(bc.ctx, "error", "Sub task already in process for account "+accId)
		bc.accountsMu.Unlock()
		return
	}
	bc.accountsMu.Unlock()

	address := bc.getAddressFromProfileId(addressProfileId)
	card := bc.getCardFromProfileId(cardProfileId)

	task := Task{
		amzTask: request.AmazonTask{
			Address: request.Address{
				FullName:    address.Name,
				Line1:       address.Line1,
				Line2:       address.Line2,
				City:        address.City,
				State:       address.State,
				Zip:         address.Zip,
				PhoneNumber: address.Phone,
			},
			Payment: request.Payment{
				CardHolder: card.Name,
				CardNumber: card.Number,
				ExpMonth:   card.ExpMonth,
				ExpYear:    card.ExpYear,
				CVV:        card.CVV,
			},
			Cookies: cookies,
		},
	}

	parsedProxy, err := forgeShr.ParseProxyString(proxy)
	if err == nil {
		task.amzTask.UserInfo.Proxy = parsedProxy
	}

	updateStatus := func(status string) {
		// return if task stopped
		if task.stopped == nil {
			return
		}

		// emit status event
		bc.accountsMu.Lock()
		for i, _ := range bc.accounts {
			newAcc := bc.accounts[i]

			if newAcc.Id == accId {
				newAcc.Status = status
				bc.accounts[i] = newAcc
			}
		}
		runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
		bc.accountsMu.Unlock()
	}

	task.amzTask.UpdateStatus = updateStatus

	go func(id string) {
		_, err := task.StartInfo()
		bc.accountsMu.Lock()
		delete(bc.accTasks, id)
		bc.accountsMu.Unlock()
		if err != nil {
			bc.accountsMu.Lock()
			for i, _ := range bc.accounts {
				newAcc := bc.accounts[i]

				if newAcc.Id == id {
					newAcc.Status = "Error: " + err.Error()
					bc.accounts[i] = newAcc
				}
			}

			runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
			bc.accountsMu.Unlock()
		}

		// if success {
		// 	bc.accountsMu.Lock()
		// 	for i, _ := range bc.accounts {
		// 		newAcc := bc.accounts[i]

		// 		if newAcc.Id == id {
		// 			bc.accounts[i] = newAcc
		// 		}
		// 	}
		// 	runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
		// 	bc.accountsMu.Unlock()

		// 	bc.SaveAccounts(bc.accounts)
		// }
	}(accId)

	bc.accountsMu.Lock()
	bc.accTasks[accId] = task
	bc.accountsMu.Unlock()
}

// * Initializes a task and immediately starts it
func (bc *BackgroundController) CreateTask() {
	email := bc.getUnusedEmail()
	proxy := bc.getProxy()

	if bc.imapClient == nil {
		runtime.EventsEmit(bc.ctx, "error", "IMAP not logged in!")
		return
	}

	if email == "" || (proxy == "" && !bc.sc.settings.LocalHost) {
		runtime.EventsEmit(bc.ctx, "error", "No emails or proxies available")
		return
	}

	if bc.sc.settings.LocalHost {
		proxy = "Localhost"
	}

	var parsedProxy forgeShr.Proxy
	var err error

	if !bc.sc.settings.LocalHost {
		parsedProxy, err = forgeShr.ParseProxyString(proxy)
		if err != nil {
			// emit err event
			runtime.EventsEmit(bc.ctx, "error", "Invalid proxy: "+proxy)
			return
		}
	}

	password := forgeShr.GeneratePassword()

	var fname, lname string
	if bc.sc.settings.NameOverride == "" {
		fname, lname = forgeShr.GenerateName()
	} else {
		nameSpl := strings.Split(bc.sc.settings.NameOverride, " ")
		if len(nameSpl) != 2 {
			fname, lname = forgeShr.GenerateName()
		} else {
			fname = nameSpl[0]
			lname = nameSpl[1]
		}
	}

	taskId := uuid.NewString()

	bc.tasksMu.Lock()
	bc.frontendTasks = append(bc.frontendTasks, TaskState{
		Id:       taskId,
		Name:     fname + " " + lname,
		Email:    email,
		Password: password,
		Proxy:    proxy,
		Status:   "Idle",
	})
	bc.tasksMu.Unlock()

	userInfo := request.UserInfo{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
		Proxy:     parsedProxy,
		ImapInfo: shr.ImapInfo{
			ImapMu:   &bc.sc.imapMu,
			Username: bc.sc.settings.ImapConfig.Username,
			Password: bc.sc.settings.ImapConfig.Password,
		},
		SmsInfo: shr.SMSInfo{
			SMSMu:      &bc.sc.tvMu,
			Username:   bc.sc.settings.SmsConfig.Username,
			ApiKey:     bc.sc.settings.SmsConfig.ApiKey,
			Provider:   bc.sc.settings.SmsConfig.Provider,
			MaxRetries: bc.sc.settings.SmsConfig.MaxTries,
		},
		CaptchaInfo: shr.CaptchaInfo{
			APIKey:     bc.sc.settings.CaptchaKey,
			MaxRetries: bc.sc.settings.CaptchaMaxTries,
		},
	}

	badFunc := func(req *http.Request) {
		os.Exit(0)
	}

	task := Task{
		amzTask: request.AmazonTask{
			UserInfo: userInfo,
			DiscordInfo: request.DiscordInfo{
				User:           bc.sc.UserInfo.DiscordUser,
				WebhookSuccess: bc.sc.settings.Webhooks.Success,
				WebhookFail:    bc.sc.settings.Webhooks.Fail,
			},
			ImapClient: bc.imapClient,
			BadFunc:    badFunc,
		},
	}

	updateStatus := func(status string) {
		// * Return if task stopped
		if task.stopped == nil {
			return
		}

		// * Emit status event to the frontend
		bc.tasksMu.Lock()
		for i := range bc.frontendTasks {
			newTask := bc.frontendTasks[i]

			if newTask.Id == taskId {
				newTask.Status = status
				bc.frontendTasks[i] = newTask
			}

		}

		runtime.EventsEmit(bc.ctx, "tasks", bc.frontendTasks)
		bc.tasksMu.Unlock()
	}

	task.amzTask.UpdateStatus = updateStatus

	go func(id string) {
		success, err := task.Start()

		if err != nil {
			bc.tasksMu.Lock()

			for i := range bc.frontendTasks {
				newTask := bc.frontendTasks[i]

				if newTask.Id == taskId {
					newTask.Status = "Error: " + err.Error()
					bc.frontendTasks[i] = newTask
				}
			}

			runtime.EventsEmit(bc.ctx, "tasks", bc.frontendTasks)

			bc.tasksMu.Unlock()

			if !success {
				bc.dataMu.Lock()
				bc.emails[email] = false
				bc.proxies[proxy] = false
				bc.dataMu.Unlock()
			}
		}

		if success {
			// marshalD, _ := json.Marshal(task.amzTask.ARD)
			account := Account{
				Name:     fname + " " + lname,
				Email:    email,
				Password: password,
				Phone:    task.amzTask.Phone,
				Proxy:    proxy,
				Key2FA:   task.amzTask.Secret,
				// ARC:      Encrypt(task.amzTask.ARC),
				// ARD:      Encrypt(string(marshalD)),
				Cookies: task.amzTask.Cookies,
				Status:  "Idle",
			}

			bc.AddAccount(account)
		}
	}(taskId)

	bc.tasksMu.Lock()
	bc.tasks[taskId] = task
	bc.tasksMu.Unlock()
}

func proxyToString(proxy forgeShr.Proxy) string {
	if proxy.Host == "" {
		return ""
	}

	pStr := proxy.Host + ":" + proxy.Port

	if proxy.User != "" {
		pStr = pStr + ":" + proxy.User + ":" + proxy.Pass
	}

	return pStr
}

func (bc *BackgroundController) deleteFrontendTask(taskId string) {
	for i, task := range bc.frontendTasks {
		if task.Id == taskId {
			bc.frontendTasks = append(bc.frontendTasks[:i], bc.frontendTasks[i+1:]...)
		}
	}
}

func (bc *BackgroundController) updateAccIdle(accId string) {
	for i, account := range bc.accounts {
		if account.Id == accId {
			account.Status = "Idle"
			bc.accounts[i] = account
		}
	}
}

func (bc *BackgroundController) StopTask(taskId string) {
	bc.tasksMu.Lock()
	defer bc.tasksMu.Unlock()

	if task, ok := bc.tasks[taskId]; ok {
		task.Stop()
		delete(bc.tasks, taskId)
		// delete(bc.frontendTasks, taskId)
		bc.deleteFrontendTask(taskId)

		emailExists := bc.emailExistsInAccounts(task.amzTask.UserInfo.Email)
		proxyExists := bc.proxyExistsInAccounts(proxyToString(task.amzTask.UserInfo.Proxy))

		if !emailExists {
			bc.emailsMu.Lock()
			bc.emails[task.amzTask.UserInfo.Email] = false
			bc.emailsMu.Unlock()
		}

		if !proxyExists && proxyToString(task.amzTask.UserInfo.Proxy) != "" {
			bc.proxiesMu.Lock()
			bc.proxies[proxyToString(task.amzTask.UserInfo.Proxy)] = false
			bc.proxiesMu.Unlock()
		}

		runtime.EventsEmit(bc.ctx, "tasks", bc.frontendTasks)
		return
	}
}

func (bc *BackgroundController) StopAccTask(accId string) {
	bc.accountsMu.Lock()
	defer bc.accountsMu.Unlock()

	if task, ok := bc.accTasks[accId]; ok {
		task.Stop()

		delete(bc.accTasks, accId)
		bc.updateAccIdle(accId)

		runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
		return
	}
}

func (bc *BackgroundController) emailExistsInAccounts(email string) bool {
	bc.accountsMu.Lock()
	defer bc.accountsMu.Unlock()

	for _, account := range bc.accounts {
		if strings.EqualFold(account.Email, email) {
			return true
		}
	}

	return false
}

func (bc *BackgroundController) proxyExistsInAccounts(proxy string) bool {
	if proxy == "" {
		return false
	}

	bc.accountsMu.Lock()
	defer bc.accountsMu.Unlock()

	for _, account := range bc.accounts {
		if strings.EqualFold(account.Proxy, proxy) {
			return true
		}
	}

	return false
}

// * Saves Emails and Proxies to the data file when user clicks save in the data tab of the ui
func (bc *BackgroundController) SaveData(emails string, proxies string) {
	// clear weird characters such as \r or \t
	emails = strings.ReplaceAll(emails, "\r", "")
	proxies = strings.ReplaceAll(proxies, "\r", "")

	emails = strings.ReplaceAll(emails, "\t", "")
	proxies = strings.ReplaceAll(proxies, "\t", "")

	bc.emailsMu.Lock()
	bc.proxiesMu.Lock()
	defer bc.emailsMu.Unlock()
	defer bc.proxiesMu.Unlock()

	emailsSpl := strings.Split(emails, "\n")
	proxiesSpl := strings.Split(proxies, "\n")

	tmpEmailMap := make(map[string]bool)
	tmpProxyMap := make(map[string]bool)

	for _, email := range emailsSpl {
		tmpEmailMap[email] = false

		if !isValidEmail(email) {
			continue
		}

		_, exists := bc.emails[email]
		if !exists {
			bc.emails[email] = false
		}
	}

	for email, used := range bc.emails {
		if _, exists := tmpEmailMap[email]; !exists && !used {
			delete(bc.emails, email)
		}
	}

	for _, proxy := range proxiesSpl {
		tmpProxyMap[proxy] = false

		if !isValidProxy(proxy) {
			continue
		}

		_, exists := bc.proxies[proxy]
		if !exists {
			bc.proxies[proxy] = false
		}
	}

	for proxy, used := range bc.proxies {
		if _, exists := tmpProxyMap[proxy]; !exists && !used {
			delete(bc.proxies, proxy)
		}
	}

	err1 := bc.SaveEmails(emails)
	err2 := bc.SaveProxies(proxies)

	if err1 == nil && err2 == nil {
		runtime.EventsEmit(bc.ctx, "success", "Successfully saved data!")
	}
}

// * ImapLogin allows logging into the IMAP provider with the given credentials
func (bc *BackgroundController) ImapLogin(imapUsername string, imapPassword string) {
	var err error
	bc.imapClient, err = forgeImap.NewEmailClient(imapUsername, imapPassword)
	if err != nil {
		runtime.EventsEmit(bc.ctx, "error", "Failed to login to imap!")
		return
	}

	runtime.EventsEmit(bc.ctx, "success", "Successfully logged in to imap!")
}

// * TestSms allows testing the SMS provider with the given credentials, it will send back the balance of the account
func (bc *BackgroundController) TestSms(provider, username, apiKey string) {
	var balance float64

	providerNice := ""

	switch provider {
	// case "Text Verified":
	// 	providerNice = "Text Verified"

	// 	tvClient, err := forgeTv.NewTVClient(username, apiKey)
	// 	if err != nil {
	// 		runtime.EventsEmit(bc.ctx, "error", "Failed to login to Text Verified!")
	// 		return
	// 	}

	// 	balance, err = tvClient.Balance()
	// 	if err != nil {
	// 		runtime.EventsEmit(bc.ctx, "error", "Failed to get balance from Text Verified!")
	// 		return
	// 	}
	case "SMS Man":
		providerNice = "SMS Man"

		smClient, err := forgeSm.NewSMSManClient(apiKey)
		if err != nil {
			runtime.EventsEmit(bc.ctx, "error", "Failed to login to SMS Man!")
			return
		}

		balance, err = smClient.Balance()
		if err != nil {
			runtime.EventsEmit(bc.ctx, "error", "Failed to get balance from SMS Man!")
			return
		}
	case "SMS Pool":
		providerNice = "SMS Pool"

		spClient, err := forgeSp.NewSMSPoolClient(apiKey)
		if err != nil {
			runtime.EventsEmit(bc.ctx, "error", "Failed to login to SMS Pool!")
			return
		}

		balance, err = spClient.Balance()
		if err != nil {
			runtime.EventsEmit(bc.ctx, "error", "Failed to get balance from SMS Pool!")
			return
		}
	case "Daisy SMS":
		providerNice = "Daisy SMS"

		dsClient, err := forgeDs.NewDaisySMSClient(apiKey)
		if err != nil {
			runtime.EventsEmit(bc.ctx, "error", "Failed to login to Daisy SMS!")
			return
		}

		balance, err = dsClient.Balance()
		if err != nil {
			runtime.EventsEmit(bc.ctx, "error", "Failed to get balance from Daisy SMS!")
			return
		}
	default:
		runtime.EventsEmit(bc.ctx, "error", "Invalid SMS provider!")
		return
	}

	runtime.EventsEmit(bc.ctx, "success", "Successfully logged in to "+providerNice+" with balance: $"+fmt.Sprintf("%.2f", balance))
}

func generateRandLabel() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	length := 8
	rand.Seed(time.Now().UnixNano())
	s := make([]rune, length)

	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return "Forge-" + string(s)
}

func (bc *BackgroundController) GenerateHMELoop() {
	if bc.sc.settings.IcloudCookies == "" {
		runtime.EventsEmit(bc.ctx, "error", "No iCloud cookies found!")
		return
	}

	u, _ := url.Parse("https://icloud.com")
	cooks := bc.iclient.HttpClient.GetCookieJar().Cookies(u)

	if len(cooks) == 0 {
		cooksSpl := strings.Split(bc.sc.settings.IcloudCookies, "; ")
		for _, cook := range cooksSpl {
			cookSpl := strings.Split(cook, "=")

			if len(cookSpl) != 2 {
				continue
			}

			cooks = append(cooks, &http.Cookie{Name: cookSpl[0], Value: cookSpl[1]})
		}

		if len(cooks) == 0 {
			runtime.EventsEmit(bc.ctx, "error", "Failed to parse iCloud cookies!")
			return
		}

		bc.iclient.HttpClient.SetCookies(u, cooks)
	}

	bc.iMu.Lock()
	if bc.generating {
		bc.iMu.Unlock()
		runtime.EventsEmit(bc.ctx, "error", "iCloud gen already running!")
		return
	}

	bc.generating = true
	bc.iMu.Unlock()

	runtime.EventsEmit(bc.ctx, "success", "Started iCloud gen")
	go func() {
		for {
			var hmes []string
			var err error
			for i := 0; i < 2; i++ {
				var hme string

				randLabel := generateRandLabel()

				hme, err = bc.iclient.ReserveHME(randLabel, "")
				if err != nil {
					continue
				}

				hmes = append(hmes, hme)
			}

			bc.dataMu.Lock()
			for _, hme := range hmes {
				bc.emails[hme] = false
			}
			bc.dataMu.Unlock()

			emailsStr := strings.Join(hmes, "\n")
			bc.SaveEmails(emailsStr)
			runtime.EventsEmit(bc.ctx, "emails", emailsStr)

			if len(hmes) > 0 {
				runtime.EventsEmit(bc.ctx, "success", "Generated HMEs: "+strings.Join(hmes, ", "))
				forgeDiscord.AlertHMESuccess(bc.sc.settings.Webhooks.Success, time.Now(), emailsStr)
			} else if err != nil {
				runtime.EventsEmit(bc.ctx, "error", "Failed to generate HMEs: "+err.Error())
			}

			time.Sleep(time.Minute * 61)
		}
	}()
}

func (bc *BackgroundController) IcloudLogin(username string, password string) {
	bc.iMu.Lock()
	defer bc.iMu.Unlock()

	if bc.iclient == nil || (bc.iclient.Username != username || bc.iclient.Password != password) {
		var err error
		bc.iclient, err = forgeIcloud.NewClient(username, password, false)
		if err != nil {
			runtime.EventsEmit(bc.ctx, "error", "Failed to init iCloud!")
			return
		}
	}

	go func() {
		err := bc.iclient.Login()
		if err != nil {
			runtime.EventsEmit(bc.ctx, "error", err)
			runtime.EventsEmit(bc.ctx, "error", "Failed to init iCloud! 2")
			return
		}
	}()

	runtime.EventsEmit(bc.ctx, "success", "Successfully started login, enter the OTP sent to your device(s)")

	otp := <-bc.otpChan
	bc.iclient.OtpChannel <- otp

	time.Sleep(time.Second * 5)
	u, _ := url.Parse("https://icloud.com")
	cooks := bc.iclient.HttpClient.GetCookieJar().Cookies(u)
	var cookies string
	for _, c := range cooks {
		cookies += c.Name + "=" + c.Value + "; "
	}
	cookies = strings.TrimSuffix(cookies, "; ")

	bc.sc.settings.IcloudConfig.Username = username
	bc.sc.settings.IcloudConfig.Password = password
	bc.sc.settings.IcloudCookies = cookies
	bc.sc.SaveSettings(bc.sc.settings, false)

	runtime.EventsEmit(bc.ctx, "success", "Successfully logged in to iCloud!")
}

func (bc *BackgroundController) SubmitOTP(otp string) {
	bc.otpChan <- otp
}
