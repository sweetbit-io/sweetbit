package updater

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseArtifactNameOutput(t *testing.T) {
	t.Parallel()

	out := parseArtifactNameOutput(`1.0.0
`)

	assert.Equal(t, out.artifactName, "1.0.0")
	assert.Equal(t, out.partitionMismatch, false)
}

func TestParseArtifactNameOutputWithMismatch(t *testing.T) {
	t.Parallel()

	out := parseArtifactNameOutput(`"time="2019-12-08T17:51:20Z" level=error msg="Mounted root '/dev/mmcblk0p3' does not match boot environment mender_boot_part: 2" module=partitions
time="2019-12-08T17:51:20Z" level=error msg="Failed to read the current active partition: No match between boot and root partitions." module=main
1.0.0
`)

	assert.Equal(t, out.artifactName, "1.0.0")
	assert.Equal(t, out.partitionMismatch, true)
}
