package updater

import (
	"bufio"
	"github.com/go-errors/errors"
	"io"
	"os/exec"
	"strings"
)

func extractUpgradeAvailable(stream io.ReadCloser) (bool, error) {
	available := false
	outputScanner := bufio.NewScanner(stream)
	for outputScanner.Scan() {
		line := strings.TrimSpace(outputScanner.Text())
		lineID := strings.SplitN(line, "=", 2)

		if lineID[0] == "upgrade_available" && lineID[1] == "1" {
			available = true
		}
	}
	err := outputScanner.Err()
	if err != nil {
		return false, errors.Errorf("unable to scan output: %v", err)
	}

	return available, nil
}

func checkUpgradeAvailable() (bool, error) {
	cmd := exec.Command("fw_printenv")

	output, err := cmd.StdoutPipe()
	if err != nil {
		return false, errors.Errorf("unable to get stdout pipe: %v", err)
	}

	defer output.Close()

	err = cmd.Start()
	if err != nil {
		return false, errors.Errorf("unable to start %s: %v", cmd.Path, err)
	}

	available, err := extractUpgradeAvailable(output)
	if err != nil {
		return false, errors.Errorf("unable to extract: %v", err)
	}

	return available, nil
}
