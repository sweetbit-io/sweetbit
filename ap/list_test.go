package ap

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestParseListNetworksPayload(t *testing.T) {
	t.Parallel()

	output := []byte("network id / ssid / bssid / flags\n0\twifi with space in name\tany\t\n1\tonion\tany\t\n2\tonion\tany\t\n")

	networks, err := parseListNetworksOutput(&output)
	if err != nil {
		log.Fatalf("Failed parsing network list output: %v", err)
	}

	if len(networks) != 3 {
		log.Fatalf("Unexpected amount of networks: %v", len(networks))
	}

	if networks[0].ssid != "wifi with space in name" {
		log.Fatalf("Unexpected first network: %v", networks[0].ssid)
	}

	if networks[1].id != "1" {
		log.Fatalf("Unexpected second network: %v", networks[0].id)
	}
}
