package antibot

import (
	"fmt"
	"math/rand"
)

var osListUa = []string{
	"17_0",
	"17_0_1",
	"17_0_2",
	"17_0_3",
	"17_1",
	"17_1_1",
	"17_1_2",
	"17_2",
	"17_2_1",
	"17_3",
	"17_3_1",
	"17_4",
	"17_4_1",
}

var osList = []string{
	"17.0",
	"17.0.1",
	"17.0.2",
	"17.0.3",
	"17.1",
	"17.1.1",
	"17.1.2",
	"17.2",
	"17.2.1",
	"17.3",
	"17.3.1",
	"17.4",
	"17.4.1",
}

func NewIOS17UserAgent() (string, string) {
	randIndex := rand.Intn(len(osListUa))
	randOs := osListUa[randIndex]
	os := osList[randIndex]

	return fmt.Sprintf("Mozilla/5.0 (iPhone; CPU iPhone OS %s like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148", randOs), os
}
