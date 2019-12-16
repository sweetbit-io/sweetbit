package updater

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func TestExtractUpgradeUnavailable(t *testing.T) {
	t.Parallel()

	stream := strings.NewReader(`stdin=serial,usbkbd
stdout=serial,vidconsole
upgrade_available=0
usb_boot=usb start; if usb dev ${devnum}; then setenv devtype usb; run scan_dev_for_boot_part; fi
vendor=raspberrypi`)

	available, err := extractUpgradeAvailable(ioutil.NopCloser(stream))

	assert.Equal(t, err, nil)
	assert.Equal(t, available, false)
}

func TestExtractUpgradeAvailable(t *testing.T) {
	t.Parallel()

	stream := strings.NewReader(`stdin=serial,usbkbd
stdout=serial,vidconsole
upgrade_available=1`)

	available, err := extractUpgradeAvailable(ioutil.NopCloser(stream))

	assert.Equal(t, err, nil)
	assert.Equal(t, available, true)
}
