package command

import (
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/v-byte-cpu/sx/pkg/scan"
)

func TestElasticCmdOptsInitCliFlags(t *testing.T) {
	t.Parallel()
	var opts elasticCmdOpts
	cmd := &cobra.Command{}

	opts.initCliFlags(cmd)
	err := cmd.ParseFlags(strings.Split(
		"--json -p 23-57,71-2733 -f ip_file.jsonl -w 300 --exit-delay 10s --timeout 2s --proto https", " "))

	require.NoError(t, err)
	require.Equal(t, true, opts.json)
	require.Equal(t, "23-57,71-2733", opts.rawPortRanges)
	require.Equal(t, "ip_file.jsonl", opts.ipFile)
	require.Equal(t, 300, opts.workers)
	require.Equal(t, 10*time.Second, opts.exitDelay)

	require.Equal(t, 2*time.Second, opts.timeout)
	require.Equal(t, "https", opts.proto)
}

func TestElasticCmdOptsParseRawOptions(t *testing.T) {
	t.Parallel()
	opts := elasticCmdOpts{
		genericScanCmdOpts: genericScanCmdOpts{
			rawPortRanges: "23-57,71-2733",
			workers:       300,
		},
		proto: "http",
	}

	err := opts.parseRawOptions()

	require.NoError(t, err)
	require.Equal(t, []*scan.PortRange{
		{StartPort: 23, EndPort: 57},
		{StartPort: 71, EndPort: 2733}}, opts.portRanges)
}
