package command

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/v-byte-cpu/sx/command/log"
	"github.com/v-byte-cpu/sx/pkg/ip"
	"github.com/v-byte-cpu/sx/pkg/scan"
	"github.com/v-byte-cpu/sx/pkg/scan/elastic"
)

var cliHTTPProtoFlag string

func init() {
	elasticCmd.Flags().StringVarP(&cliPortsFlag, "ports", "p", "", "set ports to scan")
	elasticCmd.Flags().StringVarP(&cliIPPortFileFlag, "file", "f", "", "set JSONL file with ip/port pairs to scan")
	elasticCmd.Flags().StringVar(&cliHTTPProtoFlag, "proto", "", "set protocol to use, http is used by default; only http or https are valid")
	rootCmd.AddCommand(elasticCmd)
}

var elasticCmd = &cobra.Command{
	Use: "elastic [flags] [subnet]",
	Example: strings.Join([]string{
		"elastic -p 9200 192.168.0.1/24", "elastic -p 9200-9300 10.0.0.1",
		"elastic -f ip_ports_file.jsonl", "elastic -p 9200-9300 -f ips_file.jsonl"}, "\n"),
	Short: "Perform Elasticsearch scan",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(cliHTTPProtoFlag) == 0 {
			cliHTTPProtoFlag = "http"
		}
		if cliHTTPProtoFlag != "http" && cliHTTPProtoFlag != "https" {
			return errors.New("invalid HTTP proto flag: http or https required")
		}
		if len(args) == 0 && len(cliIPPortFileFlag) == 0 {
			return errors.New("requires one ip subnet argument or file with ip/port pairs")
		}
		if len(args) == 0 {
			return
		}
		cliDstSubnet, err = ip.ParseIPNet(args[0])
		return
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		var logger log.Logger
		if logger, err = getLogger("elastic", os.Stdout); err != nil {
			return
		}

		engine := newElasticScanEngine(ctx)
		return startScanEngine(ctx, engine,
			newEngineConfig(
				withLogger(logger),
				withScanRange(&scan.Range{
					DstSubnet: cliDstSubnet,
					Ports:     cliPortRanges,
				}),
			))
	},
}

func newElasticScanEngine(ctx context.Context) scan.EngineResulter {
	// TODO custom dataTimeout
	scanner := elastic.NewScanner(cliHTTPProtoFlag, elastic.WithDataTimeout(5*time.Second))
	results := scan.NewResultChan(ctx, 1000)
	// TODO custom workerCount
	return scan.NewScanEngine(newIPPortGenerator(), scanner, results, scan.WithScanWorkerCount(50))
}