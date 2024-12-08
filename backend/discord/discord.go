package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type DiscordMsg struct {
	Username  string  `json:"username,omitempty"`
	AvatarUrl string  `json:"avatar_url,omitempty"`
	Content   string  `json:"content,omitempty"`
	Embeds    []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string     `json:"title,omitempty"`
	Url         string     `json:"url,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
	Description string     `json:"description,omitempty"`
	Color       string     `json:"color,omitempty"`
	Author      Author     `json:"author,omitempty"`
	Fields      []Field    `json:"fields,omitempty"`
	Thumbnail   Thumbnail  `json:"thumbnail,omitempty"`
	Image       Image      `json:"image,omitempty"`
	Footer      Footer     `json:"footer,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	Url     string `json:"url,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Thumbnail struct {
	Url string `json:"url,omitempty"`
}

type Image struct {
	Url string `json:"url,omitempty"`
}

type Footer struct {
	Text    string `json:"text,omitempty"`
	IconUrl string `json:"icon_url,omitempty"`
}

func AlertAccountSuccess(webhookStr string, proxy string, name string, email string, password string) error {
	t := time.Now()

	if proxy == "" {
		proxy = "N/A"
	}

	embed := Embed{
		Color: "5412454",
		Author: Author{
			Name: "The Forge Gen",
		},
		Timestamp: &t,
		Title:     "🔥 Successfully Forged Account 🔥",
		Thumbnail: Thumbnail{Url: "https://cdn.icon-icons.com/icons2/1294/PNG/512/2362134-amazon-buy-ecommerce-online-shopping_85526.png"},
		Fields: []Field{
			{
				Name:   "Site",
				Value:  "Amazon US",
				Inline: true,
			},
			{
				Name:   "Proxy",
				Value:  `||` + proxy + `||`,
				Inline: true,
			},
			{
				Name:   "Name",
				Value:  `||` + name + `||`,
				Inline: true,
			},
			{
				Name:   "Email",
				Value:  `||` + email + `||`,
				Inline: true,
			},
			{
				Name:   "Password",
				Value:  `||` + password + `||`,
				Inline: true,
			},
		},
	}

	msg := DiscordMsg{
		Username: "The Forge Bot",
		// AvatarUrl: "",
		Embeds: []Embed{embed},
	}

	return sendWebook(webhookStr, msg)
}

func AlertHMESuccess(webhookStr string, now time.Time, emails string) error {
	embed := Embed{
		Color: "5412454",
		Author: Author{
			Name: "The Forge Gen",
		},
		Timestamp: &now,
		Title:     "🔥 Successfully Forged Emails 🔥",
		Thumbnail: Thumbnail{Url: "https://cdn.freebiesupply.com/logos/large/2x/icloud-logo-svg-vector.svg"},
		Fields: []Field{
			{
				Name:   "Gen",
				Value:  "iCloud",
				Inline: true,
			},
			{
				Name:   "Emails",
				Value:  "```" + emails + "```",
				Inline: true,
			},
		},
	}

	msg := DiscordMsg{
		Username: "The Forge Bot",
		// AvatarUrl: "",
		Embeds: []Embed{embed},
	}
 
	return sendWebook(webhookStr, msg)
}

func sendWebook(url string, msg DiscordMsg) error {
	if url == "" {
		return nil
	}

	payload := new(bytes.Buffer)
	err := json.NewEncoder(payload).Encode(msg)
	if err != nil {
		return err 
	}

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		defer resp.Body.Close()

		return errors.New("Webhook failed: " + resp.Status)
	}

	return nil
}
