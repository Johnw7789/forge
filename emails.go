package main

import (
	"errors"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func isValidEmail(email string) bool {
	emailSpl := strings.Split(email, "@")

	return len(emailSpl) == 2
}

func (bc *BackgroundController) AddEmails(emails string) {
	emails = strings.TrimSuffix(emails, "\n")

	emailsSpl := strings.Split(emails, "\n")
	for _, email := range emailsSpl {
		if !isValidEmail(email) {
			emailsSpl = append(emailsSpl[:0], emailsSpl[1:]...)

			// emit err
			return
		}
	}

	for _, email := range emailsSpl {
		bc.emails[email] = false
	}

	bc.SaveEmails(emails)
}

func (bc *BackgroundController) EditEmails(emails string) {
	// trim emails so last new line is gone
	emails = strings.TrimSuffix(emails, "\n")

	emailsSpl := strings.Split(emails, "\n")
	for _, email := range emailsSpl {
		if !isValidEmail(email) {
			emailsSpl = append(emailsSpl[:0], emailsSpl[1:]...)
			// emit err
			return
		}
	}

	bc.dataMu.Lock()
	defer bc.dataMu.Unlock()

	for _, email := range emailsSpl {
		// check to see if exists in bc.emails and if it does check if in use, if not in use then delete
		if _, ok := bc.emails[email]; ok {
			if !bc.emails[email] {
				delete(bc.emails, email)
			}
		} else {
			bc.emails[email] = false
		}
	}

	bc.SaveEmails(emails)
}

func (bc *BackgroundController) LoadEmails() {
	bc.emailsMu.Lock()
	emails := ""

	emailsFile, err := GetFilePath("config", "emails.txt")
	if err != nil {
		bc.emailsMu.Unlock()
		runtime.EventsEmit(bc.ctx, "error", "Error getting emails file path")
		return
	}

	_, err = os.Stat(emailsFile)
	if os.IsNotExist(err) {
		bc.emailsMu.Unlock()
		// runtime.EventsEmit(bc.ctx, "error", "Failed to load emails")
		bc.SaveEmails(emails)
		return
	} else if err != nil {
		bc.emailsMu.Unlock()
		runtime.EventsEmit(bc.ctx, "error", "Error reading emails file")
		return
	}

	file, err := os.ReadFile(emailsFile)
	if err != nil {
		bc.emailsMu.Unlock()
		runtime.EventsEmit(bc.ctx, "error", "Error reading emails file")
		return
	}

	emails = string(file)

	emails = strings.ReplaceAll(emails, "\r", "")
	emails = strings.ReplaceAll(emails, "\t", "")

	emails = strings.TrimSuffix(emails, "\n")

	bc.emails = make(map[string]bool)
	emailsSpl := strings.Split(emails, "\n")
	for _, email := range emailsSpl {
		if !isValidEmail(email) {
			emailsSpl = append(emailsSpl[:0], emailsSpl[1:]...)
			// runtime.EventsEmit(bc.ctx, "error", fmt.Sprintf("Invalid email found in emails.txt: %s", email))
		}
	}

	for _, email := range emailsSpl {
		bc.emails[email] = false
	}

	runtime.EventsEmit(bc.ctx, "emails", emails)

	bc.emailsMu.Unlock()
}

func (bc *BackgroundController) SaveEmails(emails string) error {
	// bc.emailsMu.Lock()
	// defer bc.emailsMu.Unlock()

	// clean up any /r or /t
	emails = strings.ReplaceAll(emails, "\r", "")
	emails = strings.ReplaceAll(emails, "\t", "")

	emails = strings.TrimSuffix(emails, "\n")

	emailsFile, err := GetFilePath("config", "emails.txt")
	if err != nil {
		runtime.EventsEmit(bc.ctx, "error", "Error getting emails file path")
		return errors.New("Error getting emails file path")
	}

	err = os.WriteFile(emailsFile, []byte(emails), 0644)
	if err != nil {
		runtime.EventsEmit(bc.ctx, "error", "Error writing emails file")
		return errors.New("Error writing emails file")
	}
	return nil
}
