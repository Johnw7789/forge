package shr

import (
	crand "crypto/rand"
	"errors"
	"math/big"
	"math/rand"
	"strings"
)

var firstNames = []string{
	"Natalie", "Leah", "Zoe", "Hannah", "Audrey", "Aurora", "Bella", "Claire", "Lucy", "Stella", "Nova", "Skylar", "Zoey", "Lillian", "Aubrey", "Maya", "Madeline", "Paisley", "Adeline", "Anna", "Layla", "Alexa", "Delilah", "Samantha", "Allison", "Sarah", "Aaliyah", "Savannah", "Violet", "Autumn", "Caroline", "Ruby", "Sadie", "Ariana", "Hailey", "Kaylee", "Eleanor", "Hazel", "Genesis", "Kylie", "Eva", "Nora", "Naomi", "Piper", "Alyssa", "Brooklyn", "Taylor", "Brianna", "Katherine", "Eliza", "Maria", "Elena", "Gabriella", "Faith", "Arianna", "Melanie", "Gianna", "Isabelle", "Valentina", "Liliana", "Jade", "Willow", "Rebecca", "Clara", "Isla", "Rachel", "Amy", "Andrea", "Jasmine", "Valerie", "Adalyn", "Isabel", "Norah", "Lyla", "Michelle", "Rylee", "Charlotte", "Eden", "Emery", "Elise", "Lyla", "Aubree", "Summer", "Annabelle", "Keira", "Gracie", "Daisy", "Alana", "Molly", "Fiona", "Harmony", "Sara", "Alexis", "Sydney", "Esther", "Londyn", "Juliana", "Daniela", "Callie", "Quinn", "Hayden", "Adriana", "Teagan", "Alaina", "Angela", "Diana", "Cora", "Juliette", "Tessa", "Jayla", "Lila", "Alivia", "Presley", "Laura", "Gemma", "Lena", "Kelsey", "Gabrielle", "Camille", "Anastasia", "Jane", "Kinsley", "Lilly", "June", "April", "Maggie", "Dakota", "Rosalie", "Lia", "Paris", "Elaina", "Cecilia", "Lucia", "Brynlee", "Annie", "Holly", "Leila", "Eloise", "Maci", "Ayla", "Adelaide", "Kira", "Hope", "Elle", "Catherine", "Ruth", "Jocelyn", "Danielle", "Harley", "Lucille", "Heidi", "Evangeline", "Sienna", "Daphne", "Mckenna", "Daniella", "Demi", "Mallory", "Celeste", "Charlee", "Rylie", "Nina", "Kamila", "Vivienne", "Chelsea", "Georgia", "Talia", "Meredith", "Bianca", "Shelby", "Elsie", "Maddison", "Anya", "Josie", "Mikayla", "Alayna", "Estelle", "Phoebe", "Cassandra", "Felicity", "Harper", "Ivy", "Jordan", "Reese", "Tatum", "Veronica", "Caitlyn", "Edith", "Carmen", "Heaven", "Skyler", "Alexandria", "Selena", "Angelina", "Lana", "Miranda", "Sage", "Liana", "Lyric", "Kara", "Maeve", "Cadence", "Scarlet", "Journey", "Haven", "Nia", "Amanda", "Mariana", "Ainsley", "Athena", "Yara", "Juliet", "Jillian", "Remi", "Joanna", "Leighton", "Hanna", "Marley", "Fernanda", "Cali", "Myla", "Frida", "Leia", "Aimee", "Amara", "Evie", "Jolie", "Brooke", "Jennifer", "Margo", "Kendra", "Rowan", "Sloane", "Alma", "Lara", "Nadia", "Miriam", "Paris", "Poppy", "Cynthia", "Wren", "Sylvia", "Helen", "Tiffany", "Adelyn", "Amira", "Arielle", "Esme", "Karina", "Lorelei", "Mabel", "Sasha", "Tiana", "Emmeline", "Kori", "Elisa", "Cara", "Celia", "Hattie", "Anika", "Erika", "Lacey", "Mira", "Rosalyn", "Raegan", "Rowan", "Alina", "Bria", "Katalina", "Aviana", "Cynthia", "Elora", "Magnolia", "Alora", "Amelie", "Belle", "Elliot", "Luz", "Milani", "Raven", "Rhea", "Serena", "Virginia", "Willa", "Zara",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor", "Anderson", "Thomas", "Jackson", "White", "Harris", "Martin", "Thompson", "Garcia", "Martinez", "Robinson", "Clark", "Rodriguez", "Lewis", "Lee", "Walker", "Hall", "Allen", "Young", "Hernandez", "King", "Wright", "Lopez", "Hill", "Scott", "Green", "Adams", "Baker", "Gonzalez", "Nelson", "Carter", "Mitchell", "Perez", "Roberts", "Turner", "Phillips", "Campbell", "Parker", "Evans", "Edwards", "Collins", "Stewart", "Sanchez", "Morris", "Rogers", "Reed", "Cook", "Morgan", "Bell", "Murphy", "Bailey", "Rivera", "Cooper", "Richardson", "Cox", "Howard", "Ward", "Torres", "Peterson", "Gray", "Ramirez", "James", "Watson", "Brooks", "Kelly", "Sanders", "Price", "Bennett", "Wood", "Barnes", "Ross", "Henderson", "Coleman", "Jenkins", "Perry", "Powell", "Long", "Patterson", "Hughes", "Flores", "Washington", "Butler", "Simmons", "Foster", "Gonzales", "Bryant", "Alexander", "Russell", "Griffin", "Diaz", "Hayes",
}

func ShrBetween(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}

	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}

	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}

	return value[posFirstAdjusted:posLast]
}

func GeneratePassword() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$+?"
	length := 16

	password := make([]byte, length)

	for i := range password {
		index, err := crand.Int(crand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			panic(err)
		}
		password[i] = chars[index.Int64()]
	}

	return string(password)
}

func GenerateName() (string, string) {
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]

	return firstName, lastName
}

func GetDefaultHeaders(plain bool) map[string][]string {
	headers := map[string][]string{
		"accept-language":            {"en-US,en;q=0.9"},
		"device-memory":              {"8"},
		"Dnt":                        {"1"},
		"Downlink":                   {"5.15"},
		"Dpr":                        {"1.25"},
		"Ect":                        {"4g"},
		"Referer":                    {"https://www.amazon.com/ap/signin?openid.pape.max_auth_age=900&openid.return_to=https%3A%2F%2Fwww.amazon.com%2Fgp%2Fhomepage.html%3F_encoding%3DUTF8%26ref_%3Dnavm_accountmenu_re_signout%26path%3D%252Fgp%252Fhomepage.html%253F_encoding%253DUTF8%2526ref_%253Dnavm_accountmenu_re_signout%26useRedirectOnSuccess%3D1%26signIn%3D1%26action%3Dsign-out%26ref_%3Dnavm_accountmenu_signout&openid.assoc_handle=usflex&openid.mode=checkid_setup&openid.ns=http%3A%2F%2Fspecs.openid.net%2Fauth%2F2.0"},
		"Rtt":                        {"50"},
		"Sec-Ch-Device-Memory":       {"8"},
		"Sec-Ch-Dpr":                 {"1.25"},
		"Sec-Ch-Ua":                  {`"Google Chrome";v="112", "Chromium";v="112", "Not-A.Brand";v="24"`},
		"Sec-Ch-Ua-Mobile":           {"?1"},
		"Sec-Ch-Ua-Platform":         {`"Android"`},
		"Sec-Ch-Ua-Platform-Version": {`"6.0"`},
		"Sec-Ch-Viewport-Width":      {"339"},
		"Sec-Fetch-Dest":             {"document"},
		"Sec-Fetch-Mode":             {"navigate"},
		"Sec-Fetch-Site":             {"same-origin"},
		"Sec-Fetch-User":             {"?1"},
		"Upgrade-Insecure-Requests":  {"1"},
		"User-Agent":                 {"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Mobile Safari/537.36"},
		"Viewport-Width":             {"339"},
	}

	return headers
}

func ParseProxyString(proxyString string) (Proxy, error) {
	proxy := Proxy{}

	proxySpl := strings.Split(proxyString, ":")
	if len(proxySpl) != 2 && len(proxySpl) != 4 {
		return proxy, errors.New("Invalid proxy string")
	}

	proxy.Host = proxySpl[0]
	proxy.Port = proxySpl[1]

	if len(proxySpl) == 4 {
		proxy.User = proxySpl[2]
		proxy.Pass = proxySpl[3]
	}

	return proxy, nil
}

func GenerateRandName() (string, string) {
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]

	return firstName, lastName
}

func GenerateRandomPassword(length int) (string, error) {
	// Define the characters to be used in the password
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const specialChars = "!@$"

	// Create a buffer to store the password characters
	password := make([]byte, length-1)

	// Fill the buffer with random characters from the allowed character set
	for i := range password {
		password[i] = chars[rand.Intn(len(chars))]
	}

	// Append a random special character to the end of the password
	specialCharIndex, err := crand.Int(crand.Reader, big.NewInt(int64(len(specialChars))))
	if err != nil {
		return "", err
	}

	password = append(password, specialChars[specialCharIndex.Int64()])

	// Shuffle the password characters for additional randomness
	for i := len(password) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		password[i], password[j] = password[j], password[i]
	}

	// Convert the password buffer to a string and return
	return string(password), nil
}
