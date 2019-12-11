package upgrader

import (
	"github.com/spf13/cobra"

	"github.com/containers-ai/federatorai-operator/cmd/upgrader/influxdb"
	federatoraioperatorlog "github.com/containers-ai/federatorai-operator/pkg/log"

	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	UpgradeRootCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade ",
		Long:  "",
	}

	logOutputPath = "/var/log/alameda/federatorai-operator-upgrade.log"
)

func init() {
	UpgradeRootCmd.AddCommand(influxdb.UpgradeInfluxDBSchemaCMD)
	UpgradeRootCmd.PersistentFlags().StringVar(&logOutputPath, "log-output", logOutputPath, "File path to federatorai-operator upgrade log output")

	initLogger()
}

func initLogger() {

	cfg := federatoraioperatorlog.NewDefaultConfig()
	logger, err := federatoraioperatorlog.NewZaprLogger(cfg)
	if err != nil {
		panic(err)
	}
	log.SetLogger(logger)
}
