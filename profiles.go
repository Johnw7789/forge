package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Card struct {
	Id          string `json:"id"`
	ProfileName string `json:"profileName"`
	Name        string `json:"name"`
	Number      string `json:"number"`
	ExpMonth    string `json:"expMonth"`
	ExpYear     string `json:"expYear"`
	CVV         string `json:"cvv"`
}

type Address struct {
	Id          string `json:"id"`
	ProfileName string `json:"profileName"`
	Name        string `json:"name"`
	Line1       string `json:"line1"`
	Line2       string `json:"line2"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
	Phone       string `json:"phone"`
}

func (bc *BackgroundController) AddCards(cards []Card) {
	bc.cardsMu.Lock()
	for _, card := range cards {
		card.Id = uuid.NewString()
		bc.cards = append(bc.cards, card)
	}
	bc.cardsMu.Unlock()

	bc.SaveCards(bc.cards)
	runtime.EventsEmit(bc.ctx, "cards", bc.cards)
}

func (bc *BackgroundController) AddAddresses(addresses []Address) {
	bc.addressesMu.Lock()
	for _, address := range addresses {
		address.Id = uuid.NewString()
		bc.addresses = append(bc.addresses, address)
	}
	bc.addressesMu.Unlock()

	bc.SaveAddresses(bc.addresses)
	runtime.EventsEmit(bc.ctx, "addresses", bc.addresses)
}

func (bc *BackgroundController) AddCard(card Card) {
	bc.cardsMu.Lock()
	card.Id = uuid.NewString()
	bc.cards = append(bc.cards, card)
	bc.cardsMu.Unlock()

	bc.SaveCards(bc.cards)
	runtime.EventsEmit(bc.ctx, "cards", bc.cards)
}

func (bc *BackgroundController) AddAddress(address Address) {
	bc.addressesMu.Lock()
	address.Id = uuid.NewString()
	bc.addresses = append(bc.addresses, address)
	bc.addressesMu.Unlock()

	bc.SaveAddresses(bc.addresses)
	runtime.EventsEmit(bc.ctx, "addresses", bc.addresses)
}

func (bc *BackgroundController) EditCard(card Card) {
	bc.cardsMu.Lock()
	for i, c := range bc.cards {
		if c.Id == card.Id {
			bc.cards[i] = card
			break
		}
	}
	bc.cardsMu.Unlock()

	bc.SaveCards(bc.cards)
	runtime.EventsEmit(bc.ctx, "cards", bc.cards)
}

func (bc *BackgroundController) EditAddress(address Address) {
	bc.addressesMu.Lock()
	for i, a := range bc.addresses {
		if a.Id == address.Id {
			bc.addresses[i] = address
			break
		}
	}

	bc.addressesMu.Unlock()

	bc.SaveAddresses(bc.addresses)
	runtime.EventsEmit(bc.ctx, "addresses", bc.addresses)
}

// func (bc *BackgroundController) JigAddress(address Address) {
// 	address.Line1 = jigAddress(address.Line1)

// 	bc.EditAddress(address)
// }

func (bc *BackgroundController) DeleteCard(card Card) {
	bc.cardsMu.Lock()
	for i, c := range bc.cards {
		if c.Id == card.Id {
			bc.cards = append(bc.cards[:i], bc.cards[i+1:]...)
			break
		}
	}
	bc.cardsMu.Unlock()

	bc.SaveCards(bc.cards)
	runtime.EventsEmit(bc.ctx, "cards", bc.cards)
}

func (bc *BackgroundController) DeleteAddress(address Address) {
	bc.addressesMu.Lock()
	for i, a := range bc.addresses {
		if a.Id == address.Id {
			bc.addresses = append(bc.addresses[:i], bc.addresses[i+1:]...)
			break
		}
	}
	bc.addressesMu.Unlock()

	bc.SaveAddresses(bc.addresses)
	runtime.EventsEmit(bc.ctx, "addresses", bc.addresses)
}

func (bc *BackgroundController) LoadCards() error {
	bc.cardsMu.Lock()

	cards := []Card{}

	cardsFile, err := GetFilePath("config", "cards.json")
	if err != nil {
		bc.cardsMu.Unlock()
		return err
	}

	_, err = os.Stat(cardsFile)
	if os.IsNotExist(err) {
		bc.cardsMu.Unlock()
		bc.SaveCards(cards)
		return nil
	} else if err != nil {
		bc.cardsMu.Unlock()
		return err
	}

	file, err := os.ReadFile(cardsFile)
	if err != nil {
		bc.cardsMu.Unlock()
		return nil
	}

	err = json.Unmarshal(file, &cards)
	if err != nil {
		bc.cardsMu.Unlock()
		return nil
	}

	bc.cards = cards
	runtime.EventsEmit(bc.ctx, "cards", cards)

	bc.cardsMu.Unlock()
	return nil
}

func (bc *BackgroundController) LoadAddresses() error {
	bc.addressesMu.Lock()

	addresses := []Address{}

	addressesFile, err := GetFilePath("config", "addresses.json")
	if err != nil {
		bc.addressesMu.Unlock()
		return err
	}

	_, err = os.Stat(addressesFile)
	if os.IsNotExist(err) {
		bc.addressesMu.Unlock()
		bc.SaveAddresses(addresses)
		return nil
	}

	file, err := os.ReadFile(addressesFile)
	if err != nil {
		bc.addressesMu.Unlock()
		return nil
	}

	err = json.Unmarshal(file, &addresses)
	if err != nil {
		bc.addressesMu.Unlock()
		return nil
	}

	bc.addresses = addresses
	runtime.EventsEmit(bc.ctx, "addresses", addresses)

	bc.addressesMu.Unlock()
	return nil
}

func (bc *BackgroundController) SaveCards(cards []Card) error {
	bc.cardsMu.Lock()
	defer bc.cardsMu.Unlock()

	cardsFile, err := GetFilePath("config", "cards.json")
	if err != nil {
		return err
	}

	result, err := json.Marshal(cards)
	if err != nil {
		return errors.New("failed to save cards")
	}

	err = os.WriteFile(cardsFile, result, 0644)
	if err != nil {
		return errors.New("failed to save cards")
	}

	return nil
}

func (bc *BackgroundController) SaveAddresses(addresses []Address) error {
	bc.addressesMu.Lock()
	defer bc.addressesMu.Unlock()

	addressesFile, err := GetFilePath("config", "addresses.json")
	if err != nil {
		return err
	}

	result, err := json.Marshal(addresses)
	if err != nil {
		return errors.New("failed to save addresses")
	}

	err = os.WriteFile(addressesFile, result, 0644)
	if err != nil {
		return errors.New("failed to save addresses")
	}

	return nil
}

func (bc *BackgroundController) JigAddress(address string) string {
	// Define replacement map
	replacements := map[string]string{
		" st ":        " street ",
		" st":         " street",
		"st ":         "street ",
		" street ":    " st ",
		" street":     " st",
		" street,":    " st,",
		" street, ":   " st, ",
		" street.":    " st.",
		" street. ":   " st. ",
		" street;":    " st;",
		" ln ":        " lane ",
		" ln":         " lane",
		" lane ":      " ln ",
		" lane":       " ln",
		" lane,":      " ln,",
		" lane, ":     " ln, ",
		" lane.":      " ln.",
		" lane. ":     " ln. ",
		" lane;":      " ln;",
		" ave ":       " avenue ",
		" ave":        " avenue",
		" avenue ":    " ave ",
		" avenue":     " ave",
		" avenue,":    " ave,",
		" avenue, ":   " ave, ",
		" avenue.":    " ave.",
		" avenue. ":   " ave. ",
		"avenue ":     "ave ",
		" rd ":        " road ",
		" rd":         " road",
		" road ":      " rd ",
		" road":       " rd",
		" road,":      " rd,",
		" road, ":     " rd, ",
		" road.":      " rd.",
		" dr ":        " drive ",
		" dr":         " drive",
		" de":         "dr ",
		" drive ":     " dr ",
		" ct ":        " court ",
		" court ":     " ct ",
		" pl ":        " place ",
		" place ":     " pl ",
		" blvd ":      " boulevard ",
		" boulevard ": " blvd ",
	}

	// Convert address to lowercase for case-insensitive replacement
	addressLower := strings.ToLower(address)

	// Apply suffix switch
	for oldSuffix, newSuffix := range replacements {
		addressLower = strings.Replace(addressLower, oldSuffix, newSuffix, -1)
	}

	// Continue with other modifications
	// Define replacement map for characters
	charReplacements := map[rune]rune{
		'E': '3',
		'e': '3',
		'O': '0',
		'o': '0',
		'I': '1',
		'i': '1',
		'A': '4',
		'a': '4',
	}

	var jiggledAddress strings.Builder

	// Randomly choose whether to add characters at the start or end
	if rand.Intn(2) == 0 {
		jiggledAddress.WriteString(randomString()) // Add random string at the start
		jiggledAddress.WriteRune(' ')              // Add space
	}

	// Iterate over each character in the address
	for _, char := range addressLower {
		// Randomly choose whether to replace the character or not
		if rand.Intn(2) == 0 {
			// Check if the character needs to be replaced
			if replacement, ok := charReplacements[char]; ok {
				// Replace the character if it's in the replacements map
				jiggledAddress.WriteRune(replacement)
				continue
			}
		}
		// Randomly choose whether to switch the capitalization or not
		if rand.Intn(2) == 0 {
			// Switch the capitalization of the character
			if char >= 'a' && char <= 'z' {
				jiggledAddress.WriteRune(char - 32) // Convert lowercase to uppercase
			} else if char >= 'A' && char <= 'Z' {
				jiggledAddress.WriteRune(char + 32) // Convert uppercase to lowercase
			} else {
				jiggledAddress.WriteRune(char) // Keep the character unchanged if not a letter
			}
		} else {
			jiggledAddress.WriteRune(char) // Keep the character unchanged if not switched
		}
	}

	// Randomly choose whether to add characters at the end
	if rand.Intn(2) == 0 {
		jiggledAddress.WriteRune(' ')              // Add space
		jiggledAddress.WriteString(randomString()) // Add random string at the end
	} else {
		if rand.Intn(1) == 0 {
			jiggledAddress.WriteString(" Apt") // Add random string at the end
		}
	}

	return jiggledAddress.String()
}

// Function to generate a random string of random length between 1 and 4
func randomString() string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123"
	length := rand.Intn(3) + 2 // Random length between 1 and 4
	var result strings.Builder
	for i := 0; i < length; i++ {
		result.WriteByte(charSet[rand.Intn(len(charSet))])
	}
	return result.String()
}
