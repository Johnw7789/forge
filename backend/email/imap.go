package email

import (
	"errors"
	"io"
	"net/mail"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type EmailClient struct {
	Username string
	Password string
	Running  bool
	Client   *client.Client
	ImapMu   sync.Mutex
}

func InitImap(username, password string) (*EmailClient, error) {
	var err error
	ec, err := NewEmailClient(username, password)
	if err != nil {
		return nil, err
	}

	err = ec.imapLogin()
	if err != nil {
		return nil, errors.New("failed to login with IMAP")
	}

	return ec, nil
}

func NewEmailClient(user, pass string) (*EmailClient, error) {
	c := &EmailClient{
		Username: user,
		Password: pass,
	}

	err := c.imapLogin()
	if err != nil {
		return c, errors.New("Failed to login to IMAP")
	}

	return c, nil
}

func (e *EmailClient) imapLogin() error {
	// Connect to server
	var err error

	var domain string = "imap.mail.me.com"

	if strings.Contains(e.Username, "@gmail.com") {
		domain = "imap.gmail.com"
	} else if strings.Contains(e.Username, "@outlook.com") {
		domain = "imap-mail.outlook.com"
	} else if !strings.Contains(e.Username, "@icloud.com") {
		return errors.New("Invalid email domain")
	}

	e.Client, err = client.DialTLS(domain+":993", nil)
	if err != nil {
		return err
	}

	err = e.Client.Login(e.Username, e.Password)
	if err != nil {
		return err
	}

	return nil
}

func (e *EmailClient) searchLast10Messages(timeBefore time.Time, email string) (string, error) {
	mbox, err := e.Client.Select("Inbox", false)
	if err != nil {
		return "", err
	}

	if mbox.Messages == 0 {
		return "", errors.New("No messages in inbox")
	}

	// Get the latest message
	seqset := new(imap.SeqSet)

	size := 1

	if mbox.Messages > 10 {
		seqset.AddRange(mbox.Messages-10, mbox.Messages)
		size = 10
		// messages = make(chan *imap.Message, 10)
	} else {
		seqset.AddRange(1, mbox.Messages)
		size = int(mbox.Messages)
		// messages = make(chan *imap.Message, mbox.Messages)
	}

	messages := make(chan *imap.Message, size)

	// Get the whole message body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}

	// messages := make(chan *imap.Message, mbox.Messages)

	done := make(chan error, 1)
	go func() {
		done <- e.Client.Fetch(seqset, items, messages)
	}()

	for msg := range messages {
		if msg == nil {
			continue
		}

		r := msg.GetBody(section)
		if r == nil {
			continue
		}

		m, err := mail.ReadMessage(r)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(m.Body)
		if err != nil {
			continue
		}

		addresses, err := m.Header.AddressList("To")
		if err != nil {
			continue
		}

		if len(addresses) == 0 {
			continue
		}

		// Get the email address. so we can check if it's the right one when running a bunch of tasks at once
		address := addresses[0]
		if address == nil {
			continue
		}

		addressStr := address.Address

		timeReceived, err := m.Header.Date()
		if err != nil {
			continue
		}

		if timeBefore.Before(timeReceived) && strings.EqualFold(addressStr, email) && strings.Contains(strings.ToLower(string(body)), "amazon") && (strings.Contains(strings.ToLower(string(body)), "otp") || strings.Contains(strings.ToLower(string(body)), "verification")) {
			return string(body), nil
		}
	}

	return "", nil
}

func (e *EmailClient) FetchOtp(email string, timeBefore time.Time) (string, error) {
	var otp string
	var err error
	var msg string

	for i := 0; i < 60; i++ {
		e.ImapMu.Lock()
		msg, err = e.searchLast10Messages(timeBefore, email)
		e.ImapMu.Unlock()

		re := regexp.MustCompile(`<p>\s*(\d{6})\s*</p>`)
		matches := re.FindStringSubmatch(msg)
		if len(matches) > 1 {
			otp = matches[1] // Extract the first captured group which is the OTP
			if len(otp) == 6 {
				break
			}
		}

		re2 := regexp.MustCompile(`[<>]\b\d{6}\b[<>]`)
		matches2 := re2.FindAllString(msg, -1)

		if len(matches2) > 0 {
			otp = matches2[0]
			if len(otp) == 8 {
				otp = matches2[0][1:7]
				break
			}
		}

		split1 := strings.Split(msg, `3D"otp">`)
		if len(split1) >= 2 {
			str1 := split1[1]

			split2 := strings.Split(str1, `</p>`)

			otp = split2[0]

			if len(otp) == 6 {
				break
			}
		}

		// delay for  3s
		time.Sleep(3 * time.Second)
	}

	return otp, err
}
