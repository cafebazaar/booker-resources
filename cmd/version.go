package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/cafebazaar/booker-resources/common"
)

// serveCmd represents the serve command
var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version info of the binary",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printVersion() {
	previousLevel := logrus.GetLevel()
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Infof("Version:       %s", common.Version)
	logrus.Infof("Build Time:    %s", common.BuildTime)
	logrus.Infof("Log Level:     %s", previousLevel)
	logrus.SetLevel(previousLevel)
}
