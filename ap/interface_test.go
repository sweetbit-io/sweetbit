package ap

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestParseApInterfaceInfoOutput(t *testing.T) {
	t.Parallel()

	output := []byte("Interface uap0\n\tifindex 12\n\twdev 0xb\n\t\n\taddr b8:27:eb:6f:34:da\n\ttype AP\n\twiphy 0\n\tchannel 6 (2437 MHz), width: 20 MHz, center1: 2437 MHz\n\ttxpower 31.00 dBm")

	info, err := parseApInterfaceInfoOutput(&output)
	if err != nil {
		log.Fatalf("Failed parsing info output: %v", err)
	}

	if info.channel != 6 {
		log.Fatalf("Unexpected channel %v instead of 6", info.channel)
	}
}
