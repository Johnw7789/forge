package main

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Account struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Proxy    string `json:"proxy"`
	Key2FA   string `json:"key2fa"`
	Cookies  string `json:"cookies"`
	Prime    bool   `json:"prime"`
	Status   string `json:"status"`
}

func (bc *BackgroundController) AddAccount(account Account) {
	bc.accountsMu.Lock()
	account.Id = uuid.NewString()
	bc.accounts = append(bc.accounts, account)
	bc.accountsMu.Unlock()

	bc.SaveAccounts(bc.accounts)
	runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
}

func (bc *BackgroundController) EditAccount(account Account) {
	bc.accountsMu.Lock()
	for i, acc := range bc.accounts {
		if acc.Id == account.Id {
			bc.accounts[i] = account
			break
		}
	}
	bc.accountsMu.Unlock()

	bc.SaveAccounts(bc.accounts)
	runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
}

func (bc *BackgroundController) DeleteAccount(account Account) {
	bc.accountsMu.Lock()
	for i, acc := range bc.accounts {
		if acc.Id == account.Id {
			bc.accounts = append(bc.accounts[:i], bc.accounts[i+1:]...)
			break
		}
	}
	bc.accountsMu.Unlock()

	bc.SaveAccounts(bc.accounts)
	runtime.EventsEmit(bc.ctx, "accounts", bc.accounts)
}

func (bc *BackgroundController) LoadAccounts() error {
	bc.accountsMu.Lock()

	accounts := []Account{}

	accountsFile, err := GetFilePath("config", "accounts.json")
	if err != nil {
		bc.accountsMu.Unlock()
		return err
	}

	_, err = os.Stat(accountsFile)
	if os.IsNotExist(err) {
		bc.accountsMu.Unlock()
		bc.SaveAccounts(accounts)
		return nil
	} else if err != nil {
		bc.accountsMu.Unlock()
		return err
	}

	file, err := os.ReadFile(accountsFile)
	if err != nil {
		bc.accountsMu.Unlock()
		return nil
	}

	err = json.Unmarshal(file, &accounts)
	if err != nil {
		bc.accountsMu.Unlock()
		return nil
	}

	for _, acc := range accounts {
		acc.Status = "Idle"
	}

	bc.accounts = accounts
	runtime.EventsEmit(bc.ctx, "accounts", accounts)

	bc.accountsMu.Unlock()
	return nil
}

func (bc *BackgroundController) SaveAccounts(accounts []Account) error {
	bc.accountsMu.Lock()
	defer bc.accountsMu.Unlock()

	for _, acc := range accounts {
		acc.Status = "Idle"
	}

	accountsFile, err := GetFilePath("config", "accounts.json")
	if err != nil {
		return err
	}

	result, err := json.Marshal(accounts)
	if err != nil {
		return errors.New("failed to save accounts")
	}

	err = os.WriteFile(accountsFile, result, 0644)
	if err != nil {
		return errors.New("failed to save accounts")
	}

	return nil
}
