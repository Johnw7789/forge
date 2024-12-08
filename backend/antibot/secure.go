package antibot

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	cr "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xdg-go/pbkdf2"
)

type DeviceData struct {
	ApplicationVersion         string
	DeviceLanguage             string
	DeviceFingerprintTimestamp string
	DeviceOSVersion            string
	DeviceName                 string
	ScreenHeightPixels         string
	ThirdPartyDeviceId         string
	TimeZone                   string
	ApplicationName            string
	ScreenWidthPixels          string
	DeviceJailbroken           bool
}

type Device struct {
	DeviceData   DeviceData
	AppUserAgent string
	WebUserAgent string
	Serial       string
}

type ChallengeData struct {
	ClientId         string
	Verifier         string
	VerifierChecksum string
}

func GenerateSecureCookie(device Device) (string, error) {
	device.DeviceData.DeviceFingerprintTimestamp = fmt.Sprintf("%d", time.Now().UnixMilli())

	jsonBytes, err := json.Marshal(device)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(jsonBytes); err != nil {
		return "", err
	}

	if err := gz.Close(); err != nil {
		return "", err
	}

	compressedData := b.Bytes()
	key := pbkdf2.Key([]byte(device.Serial), []byte("AES/CBC/PKCS7Padding"), 1000, 16, sha256.New)

	iv := make([]byte, 16)
	if _, err := cr.Read(iv); err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedData := pkcs7Pad(compressedData, aes.BlockSize)
	mode := cipher.NewCBCEncrypter(block, iv)

	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)

	hmacKey := pbkdf2.Key([]byte(device.Serial), []byte("HmacSHA256"), 1000, 32, sha256.New)

	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(append(iv, ciphertext...))
	hmacData := mac.Sum(nil)[:8]

	var frc bytes.Buffer

	frc.WriteByte(0)
	frc.Write(hmacData)
	frc.Write(iv)
	frc.Write(ciphertext)

	return base64.StdEncoding.EncodeToString(frc.Bytes()), nil
}

func GenerateChallengeData(serial string) (ChallengeData, error) {
	clientId := hex.EncodeToString([]byte(serial + "#A2IVLV5VM2W81")) // Amazon alexa device id

	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return ChallengeData{}, err
	}

	verifier := strings.TrimSuffix(base64.URLEncoding.EncodeToString(verifierBytes), "=")

	h := sha256.Sum256([]byte(verifier))

	return ChallengeData{
		ClientId:         clientId,
		Verifier:         verifier,
		VerifierChecksum: strings.TrimSuffix(base64.URLEncoding.EncodeToString(h[:]), "="),
	}, nil
}

func GetRandAlexaDevice(os string) Device {
	appVer := "2.2.595606"
	osVer := "iOS/" + os
	height := "932"
	width := "430"

	return Device{
		DeviceData: DeviceData{ // now in the order that amazon has it
			ApplicationVersion:         appVer,
			DeviceLanguage:             "en-US",
			DeviceFingerprintTimestamp: fmt.Sprintf("%d", time.Now().UnixMilli()),
			DeviceOSVersion:            osVer,
			DeviceName:                 "iPhone",
			ScreenHeightPixels:         height,
			ThirdPartyDeviceId:         strings.ToUpper(uuid.New().String()),
			TimeZone:                   "-05:00",
			ApplicationName:            "Amazon Alexa",
			ScreenWidthPixels:          width,
			DeviceJailbroken:           false,
		},
		AppUserAgent: fmt.Sprintf("AmazonWebView/Amazon/%s/iOS/%s/iPhone", appVer, osVer),
		Serial:       RandDeviceSerial(),
	}
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(data, padText...)
}

func RandDeviceSerial() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 32
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}
