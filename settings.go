package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ImapConfig struct {
	UniqueTaskClient bool   `json:"uniqueTaskClient"`
	Username         string `json:"username"`
	Password         string `json:"password"`
}

type IcloudConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SmsConfig struct {
	MaxTries int    `json:"maxTries"`
	Provider string `json:"provider"`
	Username string `json:"username"`
	ApiKey   string `json:"apiKey"`
}

type Webhooks struct {
	Success string `json:"success"`
	Fail    string `json:"fail"`
}

type Settings struct {
	LicenseKey      string       `json:"licenseKey"`
	MaxTasks        int          `json:"maxTasks"`
	LimitProxyUse   bool         `json:"limitProxyUse"`
	PersistState    bool         `json:"persistState"`
	NameOverride    string       `json:"nameOverride"`
	Webhooks        Webhooks     `json:"webhooks"`
	ImapConfig      ImapConfig   `json:"imapConfig"`
	SmsConfig       SmsConfig    `json:"smsConfig"`
	CaptchaKey      string       `json:"captchaKey"`
	CaptchaMaxTries int          `json:"captchaMaxTries"`
	IcloudConfig    IcloudConfig `json:"icloudConfig"`
	IcloudCookies   string       `json:"appleCookies"`
	LocalHost       bool         `json:"localHost"`
}

type UserInfo struct {
	DiscordUser  string `json:"discordUser"`
	DiscordImage string `json:"discordImage"`
}

type SettingsController struct {
	ctx context.Context

	tvMu   sync.Mutex
	imapMu sync.Mutex

	settings   Settings
	settingsMu sync.Mutex

	UserInfo UserInfo
}

func NewSettingsController() *SettingsController {
	return &SettingsController{imapMu: sync.Mutex{}, settingsMu: sync.Mutex{}, tvMu: sync.Mutex{}}
}

func (m *SettingsController) startup(ctx context.Context) {
	m.ctx = ctx

	err := m.LoadSettings()
	if err != nil {
		os.Exit(0)
	}

	runtime.EventsOn(m.ctx, "frontend:auth", func(a ...interface{}) {
		runtime.EventsEmit(m.ctx, "backend:key", m.settings.LicenseKey)
	})

	runtime.EventsOn(m.ctx, "frontend:init", func(a ...interface{}) {
		runtime.EventsEmit(m.ctx, "settings", m.settings)
	})
}

func (m *SettingsController) LoadSettings() error {
	m.settingsMu.Lock()

	settings := Settings{}

	settingsFile, err := GetFilePath("config", "settings.json")
	if err != nil {
		m.settingsMu.Unlock()
		return err
	}

	_, err = os.Stat(settingsFile)
	if os.IsNotExist(err) {
		settings.MaxTasks = 5
		settings.SmsConfig.MaxTries = 2
		settings.SmsConfig.Provider = "SMS Man"
		settings.CaptchaMaxTries = 2
		settings.PersistState = true
		settings.LocalHost = false
		settings.ImapConfig.UniqueTaskClient = true
		settings.LimitProxyUse = true

		m.settingsMu.Unlock()
		m.SaveSettings(settings, false)
		return nil
	} else if err != nil {
		m.settingsMu.Unlock()
		return err
	}

	file, err := os.ReadFile(settingsFile)
	if err != nil {
		m.settingsMu.Unlock()
		return nil
	}

	err = json.Unmarshal(file, &settings)
	if err != nil {
		m.settingsMu.Unlock()
		return nil
	}

	m.settings = settings
	runtime.EventsEmit(m.ctx, "settings", settings)

	m.settingsMu.Unlock()
	return nil
}

func (m *SettingsController) SaveSettings(settings Settings, emitEvent bool) {
	m.settingsMu.Lock()
	defer m.settingsMu.Unlock()

	settingsFile, err := GetFilePath("config", "settings.json")
	if err != nil && emitEvent {
		runtime.EventsEmit(m.ctx, "error", "Error getting settings path!")
		return
	}

	result, err := json.Marshal(settings)
	if err != nil && emitEvent {
		runtime.EventsEmit(m.ctx, "error", "Error saving settings!")
		return
	}

	err = os.WriteFile(settingsFile, result, 0644)
	if err != nil {
		runtime.EventsEmit(m.ctx, "error", "Error saving settings file!")
		return
	}

	m.settings = settings
	runtime.EventsEmit(m.ctx, "settings", settings)

	if emitEvent {
		runtime.EventsEmit(m.ctx, "success", "Successfully saved settings!")
	}
}

func GetFilePath(subPath, fileName string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(configDir, "ForgeApp", subPath)

	err = os.MkdirAll(filePath, 0644)
	if err != nil {
		return "", err
	}

	return filepath.Join(filePath, fileName), nil
}
